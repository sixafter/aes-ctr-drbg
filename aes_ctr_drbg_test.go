// Copyright (c) 2024-2026 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.
//
// Tests for ctrdrbg: validates AES-CTR-DRBG output, uniqueness, concurrency, async rekey, personalization.

package ctrdrbg

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test_CTRDRBG_Read ensures a new DRBG produces nonzero, apparently random output for a basic read.
func Test_CTRDRBG_Read(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	rdr, err := NewReader()
	is.NoError(err)

	buf := make([]byte, 64)
	n, err := rdr.Read(buf)
	is.NoError(err)
	is.Equal(len(buf), n)

	allZeros := true
	for _, b := range buf {
		if b != 0 {
			allZeros = false
			break
		}
	}
	is.False(allZeros, "Buffer should not be all zeros")
}

// Test_CTRDRBG_ReadZeroBytes verifies that reading a zero-length buffer returns immediately with no error.
func Test_CTRDRBG_ReadZeroBytes(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	rdr, err := NewReader()
	is.NoError(err)

	buf := make([]byte, 0)
	n, err := rdr.Read(buf)
	is.NoError(err)
	is.Equal(0, n)
}

// Test_CTRDRBG_ReadMultipleTimes checks that two consecutive Read calls from the same DRBG produce different outputs.
func Test_CTRDRBG_ReadMultipleTimes(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	rdr, err := NewReader()
	is.NoError(err)

	buf1 := make([]byte, 32)
	n, err := rdr.Read(buf1)
	is.NoError(err)
	is.Equal(len(buf1), n)

	buf2 := make([]byte, 32)
	n, err = rdr.Read(buf2)
	is.NoError(err)
	is.Equal(len(buf2), n)

	is.False(bytes.Equal(buf1, buf2), "Consecutive reads should differ")
}

// Test_CTRDRBG_ReadWithDifferentBufferSizes checks that Read works correctly for a variety of buffer sizes.
func Test_CTRDRBG_ReadWithDifferentBufferSizes(t *testing.T) {
	t.Parallel()

	sizes := []int{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048}
	for _, size := range sizes {
		size := size
		t.Run("Size_"+strconv.Itoa(size), func(t *testing.T) {
			t.Parallel()
			is := assert.New(t)

			rdr, err := NewReader()
			is.NoError(err)

			buf := make([]byte, size)
			n, err := rdr.Read(buf)
			is.NoError(err)
			is.Equal(size, n)

			// An all-zero result is a plausible random outcome for small buffers (1 in 256 for one byte),
			// so only check buffers of at least 8 bytes.
			if size >= 8 {
				allZeros := true
				for _, b := range buf {
					if b != 0 {
						allZeros = false
						break
					}
				}
				is.False(allZeros, "Buffer of size %d should not be all zeros", size)
			}
		})
	}
}

// Test_CTRDRBG_Concurrency checks that the DRBG is thread-safe and outputs are unique under high concurrency.
func Test_CTRDRBG_Concurrency(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	const numGoroutines = 100
	const bufferSize = 64

	rdr, err := NewReader()
	is.NoError(err)

	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	errCh := make(chan error, numGoroutines)
	buffers := make([][]byte, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(i int) {
			defer wg.Done()
			buf := make([]byte, bufferSize)
			if _, err := rdr.Read(buf); err != nil {
				errCh <- err
				return
			}
			buffers[i] = buf
		}(i)
	}
	wg.Wait()
	close(errCh)
	for err := range errCh {
		is.NoError(err, "Concurrent Read should not error")
	}

	// Optional uniqueness check: at least two buffers should differ
	unique := false
outer:
	for i := 0; i < numGoroutines; i++ {
		for j := i + 1; j < numGoroutines; j++ {
			if !bytes.Equal(buffers[i], buffers[j]) {
				unique = true
				break outer
			}
		}
	}
	is.True(unique, "At least two buffers should differ")
}

// Test_CTRDRBG_Stream verifies that large sequential reads (1 MiB) produce random, nonzero output.
// Reads are performed in chunks of 64 KB to comply with NIST SP 800-90A max_number_of_bits_per_request.
func Test_CTRDRBG_Stream(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	rdr, err := NewReader()
	is.NoError(err)

	const total = 1 << 20     // 1 MiB
	const chunkSize = 1 << 16 // 64 KB (NIST max)
	buf := make([]byte, total)

	// Read in 64 KB chunks
	for offset := 0; offset < total; offset += chunkSize {
		end := offset + chunkSize
		if end > total {
			end = total
		}
		n, err := rdr.Read(buf[offset:end])
		is.NoError(err)
		is.Equal(end-offset, n)
	}

	allZeros := true
	for _, b := range buf {
		if b != 0 {
			allZeros = false
			break
		}
	}
	is.False(allZeros, "Stream buffer should not be all zeros")
}

// Test_CTRDRBG_ReadAll checks that large reads (10 KiB) fill the buffer with unique, nonzero data.
func Test_CTRDRBG_ReadAll(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	rdr, err := NewReader()
	is.NoError(err)

	buf := make([]byte, 10*1024) // 10 KiB
	n, err := rdr.Read(buf)
	is.NoError(err)
	is.Equal(len(buf), n)

	allZeros := true
	for _, b := range buf {
		if b != 0 {
			allZeros = false
			break
		}
	}
	is.False(allZeros, "ReadAll buffer should not be all zeros")
}

// Test_CTRDRBG_ReadConsistency performs many sequential reads, ensuring output is always nonzero and not repeated.
func Test_CTRDRBG_ReadConsistency(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	const numReads = 50
	const bufferSize = 128

	rdr, err := NewReader()
	is.NoError(err)

	buffers := make([][]byte, numReads)
	for i := 0; i < numReads; i++ {
		buf := make([]byte, bufferSize)
		n, err := rdr.Read(buf)
		is.NoError(err)
		is.Equal(bufferSize, n)

		allZeros := true
		for _, b := range buf {
			if b != 0 {
				allZeros = false
				break
			}
		}
		is.False(allZeros, "Buffer %d should not be all zeros", i)
		buffers[i] = buf
	}
	// Ensure at least two reads differ
	unique := false
outer:
	for i := 0; i < numReads; i++ {
		for j := i + 1; j < numReads; j++ {
			if !bytes.Equal(buffers[i], buffers[j]) {
				unique = true
				break outer
			}
		}
	}
	is.True(unique, "At least two buffers should differ")
}

// Test_CTRDRBG_AsyncRekey tests that async key rotation occurs after a configured threshold and the key changes.
func Test_CTRDRBG_AsyncRekey(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	cfg := DefaultConfig()
	cfg.MaxBytesPerKey = 64                  // Small threshold to trigger rekey
	cfg.RekeyBackoff = 10 * time.Millisecond // Speed up test
	cfg.MaxRekeyAttempts = 3
	cfg.MaxInitRetries = 3
	cfg.EnableKeyRotation = true

	// DRBG instance
	d, err := newDRBG(&cfg)
	is.NoError(err)

	// Get a snapshot of the initial state pointer (or block pointer)
	initialState := d.state.Load()
	initialBlock := initialState.block

	buf := make([]byte, 128) // Exceeds MaxBytesPerKey, triggers rekey
	_, err = d.Read(buf)
	is.NoError(err)

	// Wait for async rekey to finish
	wait := time.NewTimer(500 * time.Millisecond)
	tick := time.NewTicker(10 * time.Millisecond)
	defer wait.Stop()
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			currState := d.state.Load()
			// Compare block pointers (address equality) to detect a swap
			if currState.block != initialBlock && atomic.LoadUint64(&d.usage) == 0 {
				return // success
			}
		case <-wait.C:
			t.Fatal("Timed out waiting for asyncRekey to complete")
		}
	}
}

// Test_CTRDRBG_Personalization_Changes_Stream ensures different personalization strings yield unique output streams.
func Test_CTRDRBG_Personalization_Changes_Stream(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	r1, err := NewReader(WithPersonalization([]byte("foo")))
	is.NoError(err)
	r2, err := NewReader(WithPersonalization([]byte("bar")))
	is.NoError(err)

	buf1 := make([]byte, 64)
	buf2 := make([]byte, 64)

	_, err = r1.Read(buf1)
	is.NoError(err)
	_, err = r2.Read(buf2)
	is.NoError(err)

	is.False(bytes.Equal(buf1, buf2), "Personalization should affect output")
}

// Test_CTRDRBG_Read_Shards validates that a Read call only accesses a single shard in the pool.
func Test_CTRDRBG_Read_Shards(t *testing.T) {
	t.Parallel()

	// Define table of test cases with different pool (shard) counts.
	testCases := []struct {
		name       string
		shardCount int
	}{
		{"SinglePool", 1},
		{"TwoPools", 2},
		{"EightPools", 8},
		{"SixteenPools", 16},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			is := assert.New(t)

			// hit[i] will be set true if pool[i] is accessed
			hit := make([]bool, tc.shardCount)

			// Create sync.Pool array, each tracking access via hit[i]
			pools := make([]*sync.Pool, tc.shardCount)
			for i := 0; i < tc.shardCount; i++ {
				id := i
				pools[i] = &sync.Pool{
					New: func() any {
						// Record that this shard was used.
						hit[id] = true
						cfg := DefaultConfig()
						d, _ := newDRBG(&cfg)
						return d
					},
				}
			}

			r := &reader{
				pools: pools,
			}

			buf := make([]byte, 32)
			_, err := r.Read(buf)
			is.NoError(err)

			// Ensure exactly one shard was accessed.
			used := -1
			for i, v := range hit {
				if v {
					if used != -1 {
						t.Fatalf("multiple pools were accessed: %d and %d", used, i)
					}
					used = i
				}
			}
			is.NotEqual(-1, used, "no pool was used")
			t.Logf("Selected shard: %d (shardCount=%d)", used, tc.shardCount)
		})
	}
}

// Test_DRBG_FillBlocks_ZeroAlloc checks fillBlocks for correct output and zero heap allocations.
func Test_DRBG_FillBlocks_ZeroAlloc(t *testing.T) {
	cfg := DefaultConfig()
	d, _ := newDRBG(&cfg)
	var v [16]byte
	buf := make([]byte, KeySize256)

	st := d.state.Load()

	// Warmup and baseline output
	d.fillBlocks(buf, st, &v)
	baseline := make([]byte, len(buf))
	copy(baseline, buf)

	allocs := testing.AllocsPerRun(10000, func() {
		// Fill should mutate buffer and advance counter
		d.fillBlocks(buf, st, &v)
	})
	if allocs != 0 {
		t.Fatalf("unexpected allocations in fillBlocks: %v", allocs)
	}
	// Functional: Check that the buffer is not all zero
	allZero := true
	for _, b := range buf {
		if b != 0 {
			allZero = false
			break
		}
	}
	if allZero {
		t.Fatal("fillBlocks wrote only zeros")
	}
	// Functional: Output must differ from baseline (counter advanced)
	if string(baseline) == string(buf) {
		t.Fatal("fillBlocks output did not change across calls (counter not advancing?)")
	}
}

// Test_DRBG_Read_ZeroAlloc ensures DRBG.Read returns unique, random data and performs zero heap allocations.
func Test_DRBG_Read_ZeroAlloc(t *testing.T) {
	is := assert.New(t)

	cfg := DefaultConfig()
	d, _ := newDRBG(&cfg)
	buf := make([]byte, 32)

	// Warm up, baseline
	_, err := d.Read(buf)
	is.NoError(err)
	baseline := make([]byte, len(buf))
	copy(baseline, buf)

	allocs := testing.AllocsPerRun(10000, func() {
		_, err = d.Read(buf)
		is.NoError(err, "Read should not error")
	})
	if allocs != 0 {
		t.Fatalf("unexpected allocations: %v (expected 0)", allocs)
	}
	// Buffer filled?
	allZero := true
	for _, b := range buf {
		if b != 0 {
			allZero = false
			break
		}
	}
	if allZero {
		t.Fatal("Read wrote only zeros")
	}
	// Output differs from baseline?
	if string(baseline) == string(buf) {
		t.Fatal("Read output did not change across calls (counter not advancing?)")
	}
}

// Test_DRBG_Reader_Config verifies that a Reader reflects its config values as set via functional options.
func Test_DRBG_Reader_Config(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	want := Config{
		KeySize:                KeySize256,
		MaxBytesPerKey:         1024 * 1024,
		MaxInitRetries:         5,
		MaxRekeyAttempts:       7,
		MaxRekeyBackoff:        100 * time.Millisecond,
		RekeyBackoff:           10 * time.Millisecond,
		EnableKeyRotation:      true,
		Personalization:        []byte("reader-domain"),
		UseZeroBuffer:          true,
		DefaultBufferSize:      128,
		Shards:                 3,
		BacktrackingResistance: true,
	}

	rdr, err := NewReader(
		WithKeySize(want.KeySize),
		WithMaxBytesPerKey(want.MaxBytesPerKey),
		WithMaxInitRetries(want.MaxInitRetries),
		WithMaxRekeyAttempts(want.MaxRekeyAttempts),
		WithMaxRekeyBackoff(want.MaxRekeyBackoff),
		WithRekeyBackoff(want.RekeyBackoff),
		WithEnableKeyRotation(want.EnableKeyRotation),
		WithPersonalization(want.Personalization),
		WithUseZeroBuffer(want.UseZeroBuffer),
		WithDefaultBufferSize(want.DefaultBufferSize),
		WithShards(want.Shards),
		WithBacktrackingResistance(want.BacktrackingResistance),
	)
	is.NoError(err)

	got := rdr.Config()
	is.Equal(want.KeySize, got.KeySize)
	is.Equal(want.MaxBytesPerKey, got.MaxBytesPerKey)
	is.Equal(want.MaxInitRetries, got.MaxInitRetries)
	is.Equal(want.MaxRekeyAttempts, got.MaxRekeyAttempts)
	is.Equal(want.MaxRekeyBackoff, got.MaxRekeyBackoff)
	is.Equal(want.RekeyBackoff, got.RekeyBackoff)
	is.Equal(want.EnableKeyRotation, got.EnableKeyRotation)
	is.True(bytes.Equal(got.Personalization, want.Personalization), "Personalization does not match")
	is.Equal(want.UseZeroBuffer, got.UseZeroBuffer)
	is.Equal(want.DefaultBufferSize, got.DefaultBufferSize)
	is.Equal(want.Shards, got.Shards)
	is.Equal(want.BacktrackingResistance, got.BacktrackingResistance)
}

// Test_DRBG_CounterOverflow simulates 128-bit counter overflow and checks for correct counter wraparound.
func Test_DRBG_CounterOverflow(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	cfg := DefaultConfig()
	d, err := newDRBG(&cfg)
	is.NoError(err)

	// Set the DRBG's working counter to all 0xff (max 128-bit value).
	for i := 0; i < len(d.v); i++ {
		d.v[i] = 0xff
	}

	// Prepare output buffer (block size).
	blockSize := 16 // AES block size
	buf := make([]byte, blockSize)

	// Read a block -- should increment counter and wrap it to zero.
	_, err = d.Read(buf)
	is.NoError(err)

	// After increment, counter should be zero
	expected := make([]byte, 16)
	is.Equal(expected, d.v[:], "Counter should wrap to zero after overflow")

	// Optionally, check that output is nonzero
	allZeros := true
	for _, b := range buf {
		if b != 0 {
			allZeros = false
			break
		}
	}
	is.False(allZeros, "Output block should not be all zeros")
}

// Test_DRBG_ReseedInterval ensures that the DRBG reseeds itself after the configured interval.
func Test_DRBG_ReseedInterval(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	cfg := DefaultConfig()
	cfg.ReseedInterval = 50 * time.Millisecond
	d, err := newDRBG(&cfg)
	is.NoError(err)

	buf1 := make([]byte, 32)
	_, err = d.Read(buf1)
	is.NoError(err)
	initialTime := d.lastReseedTime

	// Wait just past the reseed interval
	time.Sleep(55 * time.Millisecond)

	buf2 := make([]byte, 32)
	_, err = d.Read(buf2)
	is.NoError(err)
	afterTime := d.lastReseedTime

	// Should have reseeded (time changed, outputs should differ)
	is.True(afterTime.After(initialTime) || afterTime != initialTime, "Should have reseeded on interval")
	is.False(bytes.Equal(buf1, buf2), "Outputs before/after reseed should differ")
}

// Test_DRBG_ReseedInterval_NoReseedOnFirstRead ensures a newly constructed instance does not reseed on its
// first output request merely because ReseedInterval is set; the interval is measured from instantiation.
func Test_DRBG_ReseedInterval_NoReseedOnFirstRead(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	cfg := DefaultConfig()
	cfg.ReseedInterval = time.Hour
	d, err := newDRBG(&cfg)
	is.NoError(err)

	initialState := d.state.Load()

	buf := make([]byte, 32)
	_, err = d.Read(buf)
	is.NoError(err)

	is.Same(initialState, d.state.Load(), "first Read should not reseed before ReseedInterval has elapsed")
}

// Test_DRBG_Reseed_RequestLimit checks that the DRBG reseeds after a set number of output requests.
func Test_DRBG_Reseed_RequestLimit(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	cfg := DefaultConfig()
	cfg.ReseedRequests = 3 // Reseed every 3 requests
	d, err := newDRBG(&cfg)
	is.NoError(err)

	blockSize := 32
	buf := make([]byte, blockSize)

	// Get baseline output
	_, err = d.Read(buf)
	is.NoError(err)
	out1 := append([]byte(nil), buf...)

	// Do enough reads to cross reseed threshold
	for i := 0; i < int(cfg.ReseedRequests); i++ {
		_, err = d.Read(buf)
		is.NoError(err)
	}

	// Should have triggered a reseed internally on this call
	_, err = d.Read(buf)
	is.NoError(err)
	out2 := append([]byte(nil), buf...)

	// Outputs before and after reseed should differ (with overwhelming probability)
	is.False(bytes.Equal(out1, out2), "Output after reseed should differ from before")
}

// Test_DRBG_ForkDetectionInterval_Config checks that ForkDetectionInterval is settable and present in config.
func Test_DRBG_ForkDetectionInterval_Config(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	cfg := DefaultConfig()
	cfg.ForkDetectionInterval = 17
	d, err := newDRBG(&cfg)
	is.NoError(err)

	is.Equal(uint64(17), d.config.ForkDetectionInterval, "ForkDetectionInterval should match config value")

	// Also test via functional option path.
	rdr, err := NewReader(WithForkDetectionInterval(42))
	is.NoError(err)
	got := rdr.Config()
	is.Equal(uint64(42), got.ForkDetectionInterval, "ForkDetectionInterval should be set via option")
}

// Test_CTRDRBG_InitShardPools_Failure ensures that initShardPools returns an error
// instead of panicking when all DRBG initialization attempts fail.
// It uses an invalid KeySize to deterministically force newDRBG() to fail.
func Test_CTRDRBG_InitShardPools_Failure(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	cfg := DefaultConfig()
	cfg.Shards = 2
	cfg.MaxInitRetries = 2

	// Force deterministic failure: invalid AES key size (AES supports 16, 24, 32 only)
	cfg.KeySize = 15

	pools, err := initShardPools(cfg)
	is.Error(err, "initShardPools should return error on initialization failure")
	is.Nil(pools, "pools should be nil when initialization fails")
}

// Test_DRBG_Read_MaxRequestSize verifies that Read rejects requests exceeding 64 KB.
func Test_DRBG_Read_MaxRequestSize(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	rdr, err := NewReader()
	is.NoError(err)

	// Exactly 64 KB should succeed
	buf64KB := make([]byte, 1<<16)
	n, err := rdr.Read(buf64KB)
	is.NoError(err, "Read of exactly 64 KB should succeed")
	is.Equal(1<<16, n, "Should read exactly 64 KB")

	// 64 KB + 1 byte should fail
	bufTooLarge := make([]byte, (1<<16)+1)
	n, err = rdr.Read(bufTooLarge)
	is.Error(err, "Read exceeding 64 KB should fail")
	is.Equal(ErrRequestTooLarge, err, "Should return ErrRequestTooLarge")
	is.Equal(0, n, "Should return 0 bytes on error")
}

// Test_DRBG_ReadWithAdditionalInput_MaxRequestSize verifies that ReadWithAdditionalInput
// rejects requests exceeding 64 KB.
func Test_DRBG_ReadWithAdditionalInput_MaxRequestSize(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	cfg := DefaultConfig()
	d, err := newDRBG(&cfg)
	is.NoError(err)

	// Exactly 64 KB should succeed
	buf64KB := make([]byte, 1<<16)
	n, err := d.ReadWithAdditionalInput(buf64KB, nil)
	is.NoError(err, "ReadWithAdditionalInput of exactly 64 KB should succeed")
	is.Equal(1<<16, n, "Should read exactly 64 KB")

	// 64 KB + 1 byte should fail
	bufTooLarge := make([]byte, (1<<16)+1)
	n, err = d.ReadWithAdditionalInput(bufTooLarge, nil)
	is.Error(err, "ReadWithAdditionalInput exceeding 64 KB should fail")
	is.Equal(ErrRequestTooLarge, err, "Should return ErrRequestTooLarge")
	is.Equal(0, n, "Should return 0 bytes on error")
}

// Test_AsyncRekey_WithZeroization_Enabled verifies that old key material is zeroized during rekey.
func Test_AsyncRekey_WithZeroization_Enabled(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	cfg := DefaultConfig()
	cfg.MaxBytesPerKey = 64
	cfg.RekeyBackoff = 10 * time.Millisecond
	cfg.MaxRekeyAttempts = 3
	cfg.MaxInitRetries = 3
	cfg.EnableKeyRotation = true
	cfg.EnableZeroization = true

	d, err := newDRBG(&cfg)
	is.NoError(err)

	// Capture the initial state
	initialState := d.state.Load()
	is.NotNil(initialState)

	// Trigger rekey by exceeding MaxBytesPerKey
	buf := make([]byte, 128)
	_, err = d.Read(buf)
	is.NoError(err)

	// Wait for async rekey to complete
	wait := time.NewTimer(500 * time.Millisecond)
	tick := time.NewTicker(10 * time.Millisecond)
	defer wait.Stop()
	defer tick.Stop()

	rekeyCompleted := false
	for !rekeyCompleted {
		select {
		case <-tick.C:
			currState := d.state.Load()
			if currState != initialState && atomic.LoadUint64(&d.usage) == 0 {
				rekeyCompleted = true
			}
		case <-wait.C:
			t.Fatal("Timed out waiting for asyncRekey to complete")
		}
	}

	// Verify old state key was zeroized
	allZero := true
	for _, b := range initialState.key {
		if b != 0 {
			allZero = false
			break
		}
	}
	is.True(allZero, "Old key should be zeroized after rekey")

	// Verify old state v was zeroized
	allZero = true
	for _, b := range initialState.v {
		if b != 0 {
			allZero = false
			break
		}
	}
	is.True(allZero, "Old state.v should be zeroized after rekey")
}

// Test_AsyncRekey_WithZeroization_Disabled verifies that old key material is NOT zeroized when disabled.
func Test_AsyncRekey_WithZeroization_Disabled(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	cfg := DefaultConfig()
	cfg.MaxBytesPerKey = 64
	cfg.RekeyBackoff = 10 * time.Millisecond
	cfg.MaxRekeyAttempts = 3
	cfg.MaxInitRetries = 3
	cfg.EnableKeyRotation = true
	cfg.EnableZeroization = false

	d, err := newDRBG(&cfg)
	is.NoError(err)

	// Capture the initial state
	initialState := d.state.Load()
	is.NotNil(initialState)

	// Save a copy of the original key to compare after rekey
	originalKey := initialState.key

	// Trigger rekey by exceeding MaxBytesPerKey
	buf := make([]byte, 128)
	_, err = d.Read(buf)
	is.NoError(err)

	// Wait for async rekey to complete
	wait := time.NewTimer(500 * time.Millisecond)
	tick := time.NewTicker(10 * time.Millisecond)
	defer wait.Stop()
	defer tick.Stop()

	rekeyCompleted := false
	for !rekeyCompleted {
		select {
		case <-tick.C:
			currState := d.state.Load()
			if currState != initialState && atomic.LoadUint64(&d.usage) == 0 {
				rekeyCompleted = true
			}
		case <-wait.C:
			t.Fatal("Timed out waiting for asyncRekey to complete")
		}
	}

	// Verify old state key was NOT zeroized (should still have original value)
	is.Equal(originalKey, initialState.key, "Old key should NOT be zeroized when EnableZeroization is false")
}

// Test_NewReader_WithZeroization_Enabled verifies that NewReader accepts WithZeroization option.
func Test_NewReader_WithZeroization_Enabled(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	rdr, err := NewReader(WithZeroization(true))
	is.NoError(err, "NewReader should succeed with zeroization enabled")
	is.NotNil(rdr, "Reader should be returned")

	cfg := rdr.Config()
	is.True(cfg.EnableZeroization, "EnableZeroization should be true")
}

// Test_NewReader_WithZeroization_Disabled verifies that NewReader accepts WithZeroization option.
func Test_NewReader_WithZeroization_Disabled(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	rdr, err := NewReader(WithZeroization(false))
	is.NoError(err, "NewReader should succeed with zeroization disabled")
	is.NotNil(rdr, "Reader should be returned")

	cfg := rdr.Config()
	is.False(cfg.EnableZeroization, "EnableZeroization should be false")
}

// Test_NewReader_Zeroization_DefaultDisabled verifies that zeroization is disabled by default.
func Test_NewReader_Zeroization_DefaultDisabled(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	rdr, err := NewReader()
	is.NoError(err, "NewReader should succeed with defaults")

	cfg := rdr.Config()
	is.False(cfg.EnableZeroization, "EnableZeroization should default to false")
}

// Test_CTRDRBG_Read_UseZeroBuffer_PartialBlock ensures UseZeroBuffer mode handles request lengths
// that are not a multiple of the AES block size without panicking.
func Test_CTRDRBG_Read_UseZeroBuffer_PartialBlock(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	rdr, err := NewReader(WithUseZeroBuffer(true))
	is.NoError(err)

	for _, size := range []int{1, 15, 17, 20, 33} {
		buf := make([]byte, size)
		is.NotPanics(func() {
			n, err := rdr.Read(buf)
			is.NoError(err)
			is.Equal(size, n)
		}, "Read of %d bytes should not panic", size)
		// Skip the all-zero check for tiny reads, where an all-zero result is a plausible random outcome.
		if size >= 8 {
			is.False(bytes.Equal(buf, make([]byte, size)), "Read of %d bytes should not be all zeros", size)
		}
	}
}

// Test_DRBG_BacktrackingResistance_MatchesSpec verifies the NIST SP 800-90A §10.2.1.5.1 step 6 update against an
// independent computation. With ctr_len = blocklen, a request's output (rounded up to whole blocks) and the update
// keystream form one contiguous AES-CTR stream starting at V+1, so the new Key||V must equal the seedlen bytes that
// follow the request's last output block. Each key size is checked for full and partial-block request lengths
// through Read, ReadWithAdditionalInput (nil additional input), and the UseZeroBuffer output path.
func Test_DRBG_BacktrackingResistance_MatchesSpec(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	paths := []struct {
		name       string
		zeroBuffer bool
		read       func(d *drbg, b []byte) (int, error)
	}{
		{"Read", false, func(d *drbg, b []byte) (int, error) { return d.Read(b) }},
		{"ReadWithAdditionalInput", false, func(d *drbg, b []byte) (int, error) { return d.ReadWithAdditionalInput(b, nil) }},
		{"Read_UseZeroBuffer", true, func(d *drbg, b []byte) (int, error) { return d.Read(b) }},
	}

	for _, p := range paths {
		for _, ks := range []KeySize{KeySize128, KeySize192, KeySize256} {
			for _, n := range []int{1, 16, 17, 48} {
				cfg := DefaultConfig()
				cfg.KeySize = ks
				cfg.UseZeroBuffer = p.zeroBuffer
				cfg.BacktrackingResistance = true
				d, err := newDRBG(&cfg)
				is.NoError(err)

				key0 := append([]byte(nil), d.state.Load().key[:ks]...)
				iv := d.v
				incV(&iv)

				out := make([]byte, n)
				_, err = p.read(d, out)
				is.NoError(err)

				// The output consumes whole blocks; the update keystream starts at the next block.
				outBlocks := ((n + 15) / 16) * 16
				seedLen := int(ks) + 16
				block, err := aes.NewCipher(key0)
				is.NoError(err)
				stream := make([]byte, outBlocks+seedLen)
				cipher.NewCTR(block, iv[:]).XORKeyStream(stream, stream)

				st := d.state.Load()
				msg := []any{"path=%s KeySize=%d n=%d", p.name, ks, n}
				is.Equal(stream[:n], out, msg...)
				is.Equal(stream[outBlocks:outBlocks+int(ks)], st.key[:ks], msg...)
				is.Equal(stream[outBlocks+int(ks):outBlocks+seedLen], st.v[:], msg...)
				is.Equal(st.v, d.v, msg...)
			}
		}
	}
}

// Test_DRBG_BacktrackingResistance_Allocs pins the allocation budget: zero when disabled, and at most
// two (the AES key schedule and the new immutable state) per request when enabled.
func Test_DRBG_BacktrackingResistance_Allocs(t *testing.T) {
	is := assert.New(t)

	for _, enabled := range []bool{false, true} {
		cfg := DefaultConfig()
		cfg.BacktrackingResistance = enabled
		d, err := newDRBG(&cfg)
		is.NoError(err)

		buf := make([]byte, 32)
		allocs := testing.AllocsPerRun(1000, func() { _, _ = d.Read(buf) })
		if enabled {
			is.LessOrEqual(allocs, 2.0, "BacktrackingResistance must allocate at most 2 per request")
		} else {
			is.Equal(0.0, allocs, "disabled BacktrackingResistance must not allocate")
		}
	}
}

// Test_CTRDRBG_BacktrackingResistance_Concurrency exercises the update concurrently with asynchronous key
// rotation, zeroization, health testing, request-count reseeds, and explicit Reseed calls.
func Test_CTRDRBG_BacktrackingResistance_Concurrency(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	rdr, err := NewReader(
		WithBacktrackingResistance(true),
		WithEnableKeyRotation(true),
		WithMaxBytesPerKey(64),
		WithRekeyBackoff(time.Millisecond),
		WithZeroization(true),
		WithContinuousHealthTest(true),
		WithReseedRequests(7),
		WithShards(2),
	)
	is.NoError(err)

	const numGoroutines = 32
	const iterations = 500
	var wg sync.WaitGroup
	errCh := make(chan error, numGoroutines)
	for g := 0; g < numGoroutines; g++ {
		wg.Add(1)
		go func(g int) {
			defer wg.Done()
			buf := make([]byte, 48)
			for i := 0; i < iterations; i++ {
				var err error
				switch {
				case g == 0 && i%50 == 0:
					err = rdr.Reseed(nil)
				case i%2 == 0:
					_, err = rdr.Read(buf)
				default:
					_, err = rdr.ReadWithAdditionalInput(buf[:17], []byte("additional-input"))
				}
				if err != nil {
					errCh <- err
					return
				}
			}
		}(g)
	}
	wg.Wait()
	close(errCh)
	for err := range errCh {
		is.NoError(err, "concurrent operation with BacktrackingResistance should not error")
	}
}

// Test_DRBG_BacktrackingResistance_CAVP verifies instantiation (§10.2.1.3.1) and generation (§10.2.1.5.1) against
// NIST CAVP CTR_DRBG known-answer vectors (drbgvectors_no_reseed/CTR_DRBG.rsp, "no df", PredictionResistance =
// False, AdditionalInputLen = 0, COUNT = 0 of each parameter set): Instantiate(EntropyInput, PersonalizationString),
// Generate(512 bits), Generate(512 bits), and compare the second output to ReturnedBits.
//
// The vectors run through the real constructor: crypto/rand.Reader is temporarily replaced with the vector's
// EntropyInput, so this test must not run in parallel with other tests; the original reader is restored on exit.
func Test_DRBG_BacktrackingResistance_CAVP(t *testing.T) {
	is := assert.New(t)

	originalReader := rand.Reader
	defer func() { rand.Reader = originalReader }()

	vectors := []struct {
		name            string
		keySize         KeySize
		entropyInput    string
		personalization string
		returnedBits    string
	}{
		{
			name:         "AES-128 no df, PersonalizationStringLen = 0",
			keySize:      KeySize128,
			entropyInput: "ce50f33da5d4c1d3d4004eb35244b7f2cd7f2e5076fbf6780a7ff634b249a5fc",
			returnedBits: "6545c0529d372443b392ceb3ae3a99a30f963eaf313280f1d1a1e87f9db373d3" +
				"61e75d18018266499cccd64d9bbb8de0185f213383080faddec46bae1f784e5a",
		},
		{
			name:            "AES-128 no df, PersonalizationStringLen = 256",
			keySize:         KeySize128,
			entropyInput:    "cee23de86a69c7ef57f6e1e12bd16e35e51624226fa19597bf93ec476a44b0f2",
			personalization: "a2ef16f226ea324f23abd59d5e3c660561c25e73638fe21c87566e86a9e04c3e",
			returnedBits: "2a76d71b329f449c98dc08fff1d205a2fbd9e4ade120c7611c225c984eac8531" +
				"288dd3049f3dc3bb3671501ab8fbf9ad49c86cce307653bd8caf29cb0cf07764",
		},
		{
			name:         "AES-192 no df, PersonalizationStringLen = 0",
			keySize:      KeySize192,
			entropyInput: "f1ef7eb311c850e189be229df7e6d68f1795aa8e21d93504e75abe78f041395873540386812a9a2a",
			returnedBits: "6bb0aa5b4b97ee83765736ad0e9068dfef0ccfc93b71c1d3425302ef7ba4635f" +
				"fc09981d262177e208a7ec90a557b6d76112d56c40893892c3034835036d7a69",
		},
		{
			name:            "AES-192 no df, PersonalizationStringLen = 320",
			keySize:         KeySize192,
			entropyInput:    "d76697495add9baac8cfdf90e0a6f3381a52b2889f840a24cfc9606f97c49446202e70583fde09dc",
			personalization: "d64e546523fb082b5e87bbe6b66a8e6a7940bb06ed17fe833c4e9bbb211396c227f83671c924240e",
			returnedBits: "5693f68ec7a8d2631dee9d371503e05cbd5fb872506f040050b1d441a94bd74e" +
				"c06f24c1cc4363ece1d8b2295e74389b6dd8b1e5438bceab2f2bc2af5a9d8da8",
		},
		{
			name:         "AES-256 no df, PersonalizationStringLen = 0",
			keySize:      KeySize256,
			entropyInput: "df5d73faa468649edda33b5cca79b0b05600419ccb7a879ddfec9db32ee494e5531b51de16a30f769262474c73bec010",
			returnedBits: "d1c07cd95af8a7f11012c84ce48bb8cb87189e99d40fccb1771c619bdf82ab22" +
				"80b1dc2f2581f39164f7ac0c510494b3a43c41b7db17514c87b107ae793e01c5",
		},
		{
			name:            "AES-256 no df, PersonalizationStringLen = 384",
			keySize:         KeySize256,
			entropyInput:    "22a89ee0e37b54ea636863d9fed10821f1952a428488d528eceb9d2ec69d573ec6216216fb3e8f72a148a5ada9d620b1",
			personalization: "953c10badcbcd45fb4e5475826477fc137ac96a49ad5005fb14bdaf6468ae7f46c5d0de22d304afc67989615adc2e983",
			returnedBits: "f7fab6a6fcf445f0a0434b2aa0c610bdef5489ecd95414634623add18a9f888b" +
				"ca6be151312d1b9e8f83bd0acad6234d3bccc11b63a40d6fbff448f67db0b91f",
		},
	}

	for _, v := range vectors {
		entropyInput, err := hex.DecodeString(v.entropyInput)
		is.NoError(err)
		personalization, err := hex.DecodeString(v.personalization)
		is.NoError(err)
		want, err := hex.DecodeString(v.returnedBits)
		is.NoError(err)

		cfg := DefaultConfig()
		cfg.KeySize = v.keySize
		cfg.Personalization = personalization
		cfg.BacktrackingResistance = true

		rand.Reader = bytes.NewReader(entropyInput)
		d, err := newDRBG(&cfg)
		rand.Reader = originalReader
		is.NoError(err, v.name)

		got := make([]byte, len(want))
		_, err = d.Read(got)
		is.NoError(err, v.name)
		_, err = d.Read(got)
		is.NoError(err, v.name)
		is.Equal(want, got, v.name)
	}
}

// Test_NewReader_BacktrackingResistance_PersonalizationLength verifies that, with BacktrackingResistance enabled,
// NewReader enforces the NIST SP 800-90A §10.2.1 Table 3 maximum personalization string length (seedlen) for each
// key size, and that the default mode continues to accept longer personalization strings.
func Test_NewReader_BacktrackingResistance_PersonalizationLength(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	for _, ks := range []KeySize{KeySize128, KeySize192, KeySize256} {
		seedLen := int(ks) + 16

		_, err := NewReader(WithKeySize(ks), WithBacktrackingResistance(true), WithPersonalization(make([]byte, seedLen)))
		is.NoError(err, "personalization of seedlen bytes must be accepted (KeySize=%d)", ks)

		_, err = NewReader(WithKeySize(ks), WithBacktrackingResistance(true), WithPersonalization(make([]byte, seedLen+1)))
		is.Equal(ErrPersonalizationTooLong, err, "personalization longer than seedlen must be rejected (KeySize=%d)", ks)

		_, err = NewReader(WithKeySize(ks), WithPersonalization(make([]byte, seedLen+1)))
		is.NoError(err, "default mode must continue to accept longer personalization (KeySize=%d)", ks)
	}
}

// Test_DRBG_BacktrackingResistance_CAVP_AdditionalInput verifies additional input handling in generation
// (§10.2.1.5.1 steps 2 and 6) against NIST CAVP CTR_DRBG known-answer vectors (drbgvectors_no_reseed/CTR_DRBG.rsp,
// "no df", PredictionResistance = False, AdditionalInputLen = seedlen, COUNT = 0 of each parameter set):
// Instantiate(EntropyInput, PersonalizationString), Generate(512 bits, AdditionalInput1),
// Generate(512 bits, AdditionalInput2), and compare the second output to ReturnedBits.
//
// The vectors run through the real constructor: crypto/rand.Reader is temporarily replaced with the vector's
// EntropyInput, so this test must not run in parallel with other tests; the original reader is restored on exit.
func Test_DRBG_BacktrackingResistance_CAVP_AdditionalInput(t *testing.T) {
	is := assert.New(t)

	originalReader := rand.Reader
	defer func() { rand.Reader = originalReader }()

	vectors := []struct {
		name             string
		keySize          KeySize
		entropyInput     string
		personalization  string
		additionalInput1 string
		additionalInput2 string
		returnedBits     string
	}{
		{
			name:             "AES-128 no df, PersonalizationStringLen = 0, AdditionalInputLen = 256",
			keySize:          KeySize128,
			entropyInput:     "6bd4f2ae649fc99350951ff0c5d460c1a9214154e7384975ee54b34b7cae0704",
			additionalInput1: "ecd4893b979ac92db1894ae3724518a2f78cf2dbe2f6bbc6fda596df87c7a4ae",
			additionalInput2: "b23e9188687c88768b26738862c4791fa52f92502e1f94bf66af017c4228a0dc",
			returnedBits: "5b2bf7a5c60d8ab6591110cbd61cd387b02de19784f496d1a109123d8b3562a5" +
				"de2dd6d5d1aef957a6c4f371cecd93c15799d82e34d6a0dba7e915a27d8e65f3",
		},
		{
			name:             "AES-128 no df, PersonalizationStringLen = 256, AdditionalInputLen = 256",
			keySize:          KeySize128,
			entropyInput:     "50b96542a1f2b8b05074051fe8fb0e45adbbd5560e3594e12d485fe1bfcb741f",
			personalization:  "820c3030f97b3ead81a93b88b871937278fd3d711d2085d9280cba394673b17e",
			additionalInput1: "1f1632058806d6d8e231288f3b15a3c324e90ccef4891bd595f09c3e80e27469",
			additionalInput2: "5cadc8bfd86d2a5d44f921f64c7d153001b9bdd7caa6618639b948ebfad5cb8a",
			returnedBits: "02b76a66f103e98d450e25e09c35337747d987471d2b3d81e03be24c7e985417" +
				"a32acd72bc0a6eddd9871410dacb921c659249b4e2b368c4ac8580fb5db559bc",
		},
		{
			name:             "AES-192 no df, PersonalizationStringLen = 0, AdditionalInputLen = 320",
			keySize:          KeySize192,
			entropyInput:     "b895b891f039052821fbb4a889fced861b96c37e36a5f4f7aa208a2bf33a896fe7e29f3f6cf0041f",
			additionalInput1: "8711df3931a9035905ebdf510b3ea3f344923b2f20a561709c0def03b9570be267e9765719a25d8a",
			additionalInput2: "03908b9fd56bad5a1645b688a4f9c2342bc0f81fed6f7449af14e92960d60127dda4e27487e77491",
			returnedBits: "1385397b1a245bd06bbb4bcf651a52b2d30867d3e2f98e9c7a9ee959d2e1bbd6" +
				"3a10054fc081cd22a7f3b65ed2f0b3a3deb389d2f336e50b359a6c3e83667fb1",
		},
		{
			name:             "AES-192 no df, PersonalizationStringLen = 320, AdditionalInputLen = 320",
			keySize:          KeySize192,
			entropyInput:     "8db47ca284294daeb303e6459422e32577ca539c00fead8839d42946bad8eb30e489c6d858763a28",
			personalization:  "c2491a6ce1ef9bd13eb66a86c0f5be3d2f3239bf7f713e830e8ae8907a2084f873ed3f5eddf5b569",
			additionalInput1: "e9a3f24180345a0c0680aacdb89ce7cf841c7ad547ef92dad54f8262445e2f0c54c8f8e42350f79c",
			additionalInput2: "e2427b93738424c0fce14cb6c5f1d6b6a0532787157b6d907bc55d1c9a670494776212a643a9fb2c",
			returnedBits: "494b1957ab382b382054f54bb62f0b5b464afb3b18492c60809c26e86e45b6b9" +
				"fa44524dd89ecdca99c60e68ed107ff336be1591cadd6fd7e35f742611809f2e",
		},
		{
			name:             "AES-256 no df, PersonalizationStringLen = 0, AdditionalInputLen = 384",
			keySize:          KeySize256,
			entropyInput:     "f45e9d040c1456f1c7f26e7f146469fbe3973007fe037239ad57623046e7ec52221b22eec208b22ac4cf4ca8d6253874",
			additionalInput1: "28819bc79b92fc8790ebdc99812cdcea5c96e6feab32801ec1851b9f46e80eb6800028e61fbccb6ccbe42b06bf5a0864",
			additionalInput2: "418ca848027e1b3c84d66717e6f31bf89684d5db94cd2d579233f716ac70ab66cc7b01a6f9ab8c7665fcc37dba4af1ad",
			returnedBits: "4f11406bd303c104243441a8f828bf0293cb20ac39392061429c3f56c1f42623" +
				"9f8f0c687b69897a2c7c8c2b4fb520b62741ffdd29f038b7c82a9d00a890a3ed",
		},
		{
			name:             "AES-256 no df, PersonalizationStringLen = 384, AdditionalInputLen = 384",
			keySize:          KeySize256,
			entropyInput:     "0dd4d80062ecc0f359efbe7723020be9b88b550fe74088094069e74428395856f63eed4f5b0e7d1e006f0eaff74f638c",
			personalization:  "d2aa2ccd4bc6537e51f6550ab6d6294547bef3e971a7f128e4436f957de9982c93ee22110b0e40ab33a7d3dfa22f599d",
			additionalInput1: "0b081bab6c74d86b4a010e2ded99d14e0c9838f7c3d69afd64f1b66377d95cdcb7f6ec5358e3516034c3339ced7e1638",
			additionalInput2: "ca818f938ae0c7f4f507e4cfec10e7baf51fe34b89a502f754d2d2be7395120fe1fb013c67ac2500b3d17b735da09a6e",
			returnedBits: "6808268b13e236f642c06deba2494496e7003c937ebf6f7cb7c92104ea090f18" +
				"484aa075560d7844a06eb559948c93b26ae40f2db98ecb53ad593eb4c78f82b1",
		},
	}

	for _, v := range vectors {
		entropyInput, err := hex.DecodeString(v.entropyInput)
		is.NoError(err)
		personalization, err := hex.DecodeString(v.personalization)
		is.NoError(err)
		additionalInput1, err := hex.DecodeString(v.additionalInput1)
		is.NoError(err)
		additionalInput2, err := hex.DecodeString(v.additionalInput2)
		is.NoError(err)
		want, err := hex.DecodeString(v.returnedBits)
		is.NoError(err)

		cfg := DefaultConfig()
		cfg.KeySize = v.keySize
		cfg.Personalization = personalization
		cfg.BacktrackingResistance = true

		rand.Reader = bytes.NewReader(entropyInput)
		d, err := newDRBG(&cfg)
		rand.Reader = originalReader
		is.NoError(err, v.name)

		got := make([]byte, len(want))
		_, err = d.ReadWithAdditionalInput(got, additionalInput1)
		is.NoError(err, v.name)
		_, err = d.ReadWithAdditionalInput(got, additionalInput2)
		is.NoError(err, v.name)
		is.Equal(want, got, v.name)
	}
}

// Test_DRBG_BacktrackingResistance_AdditionalInputSemantics verifies, with BacktrackingResistance enabled, that
// additional input longer than seedlen is rejected with ErrAdditionalInputTooLong, that shorter additional input
// is treated as zero-padded to seedlen (NIST SP 800-90A §10.2.1.5.1 step 2.2), and that empty additional input is
// treated as Null. It also verifies the default mode continues to accept longer additional input.
//
// Two instances are seeded with identical entropy, so this test temporarily replaces crypto/rand.Reader and must
// not run in parallel with other tests; the original reader is restored on exit.
func Test_DRBG_BacktrackingResistance_AdditionalInputSemantics(t *testing.T) {
	is := assert.New(t)

	originalReader := rand.Reader
	defer func() { rand.Reader = originalReader }()

	for _, ks := range []KeySize{KeySize128, KeySize192, KeySize256} {
		seedLen := int(ks) + 16
		cfg := DefaultConfig()
		cfg.KeySize = ks
		cfg.BacktrackingResistance = true

		// pair returns two instances instantiated from the same entropy.
		pair := func() (*drbg, *drbg) {
			entropy := bytes.Repeat([]byte{0x5a}, seedLen)
			rand.Reader = bytes.NewReader(entropy)
			a, err := newDRBG(&cfg)
			is.NoError(err)
			rand.Reader = bytes.NewReader(entropy)
			b, err := newDRBG(&cfg)
			is.NoError(err)
			rand.Reader = originalReader
			return a, b
		}

		buf := make([]byte, 32)
		d, _ := pair()
		_, err := d.ReadWithAdditionalInput(buf, make([]byte, seedLen))
		is.NoError(err, "additional input of seedlen bytes must be accepted (KeySize=%d)", ks)
		n, err := d.ReadWithAdditionalInput(buf, make([]byte, seedLen+1))
		is.Equal(ErrAdditionalInputTooLong, err, "additional input longer than seedlen must be rejected (KeySize=%d)", ks)
		is.Equal(0, n)

		// Short additional input (seedlen-1 bytes, ending in a partial block) must behave exactly like the
		// same input zero-padded to seedlen.
		short := bytes.Repeat([]byte{0xa5}, seedLen-1)
		padded := make([]byte, seedLen)
		copy(padded, short)
		a, b := pair()
		outA, outB := make([]byte, 48), make([]byte, 48)
		_, err = a.ReadWithAdditionalInput(outA, short)
		is.NoError(err)
		_, err = b.ReadWithAdditionalInput(outB, padded)
		is.NoError(err)
		is.Equal(outB, outA, "short additional input must be zero-padded to seedlen (KeySize=%d)", ks)

		// Empty additional input must behave exactly like nil (Null).
		a, b = pair()
		_, err = a.ReadWithAdditionalInput(outA, []byte{})
		is.NoError(err)
		_, err = b.ReadWithAdditionalInput(outB, nil)
		is.NoError(err)
		is.Equal(outB, outA, "empty additional input must be treated as Null (KeySize=%d)", ks)
	}

	r, err := NewReader()
	is.NoError(err)
	_, err = r.ReadWithAdditionalInput(make([]byte, 32), make([]byte, 64))
	is.NoError(err, "default mode must continue to accept additional input longer than seedlen")
}

// Test_DRBG_BacktrackingResistance_CAVP_PredictionResistance verifies prediction-resistance reseeding
// (NIST SP 800-90A §9.3.1 step 7 and §10.2.1.4.1) against NIST CAVP CTR_DRBG known-answer vectors
// (drbgvectors_pr_true/CTR_DRBG.rsp, "no df", PredictionResistance = True, COUNT = 0 of each listed parameter set):
// Instantiate(EntropyInput, PersonalizationString); for each of two Generate(512 bits, AdditionalInput) calls,
// Reseed(EntropyInputPR, AdditionalInput) then generate; compare the second output to ReturnedBits.
//
// The vectors run through the real constructor: crypto/rand.Reader is temporarily replaced with the vector's
// EntropyInput followed by both EntropyInputPR values, so this test must not run in parallel with other tests;
// the original reader is restored on exit.
func Test_DRBG_BacktrackingResistance_CAVP_PredictionResistance(t *testing.T) {
	is := assert.New(t)

	originalReader := rand.Reader
	defer func() { rand.Reader = originalReader }()

	vectors := []struct {
		name             string
		keySize          KeySize
		entropyInput     string
		personalization  string
		additionalInput1 string
		entropyInputPR1  string
		additionalInput2 string
		entropyInputPR2  string
		returnedBits     string
	}{
		{
			name:             "AES-128 no df, PersonalizationStringLen = 0, AdditionalInputLen = 256",
			keySize:          KeySize128,
			entropyInput:     "4d85ec7a58763c383548e70758c6db91f99b1466baa7d4f20ba1f8a44ec4d1f3",
			additionalInput1: "dd4f1c82473fe4cc56128152001547ff7234eadea5e94f0d64b0c08d24788570",
			entropyInputPR1:  "3a9e8e7aebfed99f34b16294c1da973bd706afba105d78cb132e4d3aee221cc1",
			additionalInput2: "e90046b7fcd9d86434c312b959795fbfbe666dda2b16a7b14cb0cc1d84895273",
			entropyInputPR2:  "e37f193e5f1356368aec70a0de2c330c3cf9af26c6aee46b9f114f6276ed38b4",
			returnedBits: "173a535a1ce52e3ca2b68fb2c24880a4dc2837174ce3df1c1f445f43013384a0" +
				"570ec8d5197b2d00fef3cb3556086bd82fb6075562f3adae4dca39008bcd7397",
		},
		{
			name:            "AES-128 no df, PersonalizationStringLen = 256, AdditionalInputLen = 0",
			keySize:         KeySize128,
			entropyInput:    "528ccd2d6c143800a34ad33e7f153cfaceaa2411abbaf4bfcfe9796898d0ece6",
			personalization: "07ca7a007b53f014d7f10461d6c97b5c8b0c502d3461d583f99efec69f121cd2",
			entropyInputPR1: "478fd1eaa7ed293294d370979b0f0f1948d5a3161b12eeebf2cf6bd1bf059adf",
			entropyInputPR2: "036bf1173977b323f60c6f0f603d9c50835b2afca8347dfb24e8f66604444951",
			returnedBits: "bba1ef50b4bad288897f02ac2706ee1e01488dcbe9b3d8a637921f5e788faed3" +
				"b23db63d590ceaa7607a2179192bea9aeaa85d048e7e108fee666dc646af5f0e",
		},
		{
			name:             "AES-192 no df, PersonalizationStringLen = 0, AdditionalInputLen = 320",
			keySize:          KeySize192,
			entropyInput:     "4be5afd0d9853bb3199bffb651211e4020afe4c27dd97349be9d21e260c2c4cb74014c4f9684ecab",
			additionalInput1: "22305e83b7c362f93abeae0044459b2ba9bc108526555881dc9a9e7f74afbaade44996ad0478b015",
			entropyInputPR1:  "26af31f45eed020ca9146d15d72e5e5d811a4cfe0933c5ebf2b1d1d912cd4d3374e78b24f6c7d9fd",
			additionalInput2: "587196fa6d13a85006e258cfef4bce56e916c652125fb79ffc8af34e0f2b7c357c1d61fda6f5e3fa",
			entropyInputPR2:  "bb39f25481e8de9ffd540283ac67479aae1d9cf614315db7d8ec64afa9a841ca938c92896d8e8424",
			returnedBits: "a4b1add24bd78366d77125ec100bb3c40dd1ae1a0ee8e026541462243e5d2db2" +
				"f1a47affbdd364360359dbb5e6eeb525d73a43a47531c04c03325352c6f25084",
		},
		{
			name:            "AES-192 no df, PersonalizationStringLen = 320, AdditionalInputLen = 0",
			keySize:         KeySize192,
			entropyInput:    "6e7c6d3224b574f9d197037a7770db1f005dfc82d6fd19399d73614bb398fa0e259848ded31f19e4",
			personalization: "d0a67c39d4ee83208d40debebf38332b0c09b6683a5070dfd91bdfb0aa7c60af4dc36b18926189c6",
			entropyInputPR1: "b481b84b2ab1c542e5828b3db7067729516bb2eaf6f93f65538102935cd4fd2d3675bab00980a217",
			entropyInputPR2: "9ad8c9710fbecb0ca91f0c1f07cbea29efa6753a2f39b57c9d10dc13eb2d8c39846712167b13c42e",
			returnedBits: "ce788ccd16c5c21f91d26556b7fcb0aa644a0d2692ab3cd9e8a50860d402e6b9" +
				"b8c824fe38c4c150a2ce97b939b6b633cd48991e152b8134462ea3f039b11c0b",
		},
		{
			name:             "AES-256 no df, PersonalizationStringLen = 0, AdditionalInputLen = 384",
			keySize:          KeySize256,
			entropyInput:     "c519cd4f839c9a5e936485242c4da847f47af5b544acd12b49036453655de3ed57b2ae015822a60edcfb859b00bbef2e",
			additionalInput1: "7e95d809c11a7a975230da9599d0073136bd057a1a6990f980894c01912bc39b8dac94c083933d39278f65542dd8dd4e",
			entropyInputPR1:  "0415f8fd00d4d2454706142083346304d8051b8259d08de493d991d366a42fd452d6533346896d323ba1066a485ef105",
			additionalInput2: "188357a4c4dbdd05cb6bc5d9160d5d951cef6ed9ab6c622b5a9cda2c761d9278235fbe9d6312b9d892a59abb7f3f4865",
			entropyInputPR2:  "9a146501a4dc9349b0392897263cfa414b6fc7e59d7093ab66bd265b7e41b5c2766f56d97b7d22d02220cbc389b34a6f",
			returnedBits: "cb79c79dd19a72b303a2568266abf23ad6e341cfcfee8a768af64584e4e471de" +
				"dc0d832a04e8dcf4b896ad9bf4a4a2495cc6f3efd22a68938424717dd75e83b0",
		},
		{
			name:            "AES-256 no df, PersonalizationStringLen = 384, AdditionalInputLen = 0",
			keySize:         KeySize256,
			entropyInput:    "c8f0c7b9bdf7e7d524c998aedeabb3b7dd4fa8f95c51b582010a5e09d0b4b1ad510302422df738fbef002a051543b4cc",
			personalization: "3b8f2ed03e4c4a857ebaf4730124ae5b3ffa67e294381fe085e1d5b25bb79c228f39b9912d8bf9e84cdb43fcefb552e0",
			entropyInputPR1: "c1a7b160c8e33a01fbd49743dc1161539390d9ba6b876fe63b58e8fd605b98617322578b17aca9db71e858b154f97910",
			entropyInputPR2: "d98b25eda41d15eacdecba586609d8c743704fa099ac37f9185bcb19652723d1648431a73c4399773c85caf12fcdd842",
			returnedBits: "338da59350a72bba94b217b99ff813fafa85f31fa7e991d4350d6fce39471aa2" +
				"49c6a9e7d189649ea6770a15be302b943ad403728b7387f5984c9bfb82597771",
		},
	}

	decode := func(s string) []byte {
		b, err := hex.DecodeString(s)
		is.NoError(err)
		return b
	}

	for _, v := range vectors {
		want := decode(v.returnedBits)

		cfg := DefaultConfig()
		cfg.KeySize = v.keySize
		cfg.Personalization = decode(v.personalization)
		cfg.PredictionResistance = true
		cfg.BacktrackingResistance = true

		var entropy []byte
		entropy = append(entropy, decode(v.entropyInput)...)
		entropy = append(entropy, decode(v.entropyInputPR1)...)
		entropy = append(entropy, decode(v.entropyInputPR2)...)
		rand.Reader = bytes.NewReader(entropy)

		d, err := newDRBG(&cfg)
		is.NoError(err, v.name)
		got := make([]byte, len(want))
		_, err = d.ReadWithAdditionalInput(got, decode(v.additionalInput1))
		is.NoError(err, v.name)
		_, err = d.ReadWithAdditionalInput(got, decode(v.additionalInput2))
		is.NoError(err, v.name)
		rand.Reader = originalReader

		is.Equal(want, got, v.name)

		// Without additional input, Read must take the same prediction-resistance path.
		if v.additionalInput1 == "" && v.additionalInput2 == "" {
			rand.Reader = bytes.NewReader(entropy)
			d, err = newDRBG(&cfg)
			is.NoError(err, v.name)
			_, err = d.Read(got)
			is.NoError(err, v.name)
			_, err = d.Read(got)
			is.NoError(err, v.name)
			rand.Reader = originalReader

			is.Equal(want, got, "Read: "+v.name)
		}
	}
}

// Test_DRBG_BacktrackingResistance_CAVP_Reseed verifies explicit reseeding (NIST SP 800-90A §10.2.1.4.1) against
// NIST CAVP CTR_DRBG known-answer vectors (drbgvectors_pr_false/CTR_DRBG.rsp, "no df", PredictionResistance =
// False, COUNT = 0 of each listed parameter set): Instantiate(EntropyInput, PersonalizationString),
// Reseed(EntropyInputReseed, AdditionalInputReseed), Generate(512 bits, AdditionalInput) twice, and compare the
// second output to ReturnedBits.
//
// The vectors run through the real constructor and Reseed: crypto/rand.Reader is temporarily replaced with the
// vector's EntropyInput followed by EntropyInputReseed, so this test must not run in parallel with other tests;
// the original reader is restored on exit.
func Test_DRBG_BacktrackingResistance_CAVP_Reseed(t *testing.T) {
	is := assert.New(t)

	originalReader := rand.Reader
	defer func() { rand.Reader = originalReader }()

	vectors := []struct {
		name                  string
		keySize               KeySize
		entropyInput          string
		personalization       string
		entropyInputReseed    string
		additionalInputReseed string
		additionalInput1      string
		additionalInput2      string
		returnedBits          string
	}{
		{
			name:                  "AES-128 no df, PersonalizationStringLen = 0, AdditionalInputLen = 256",
			keySize:               KeySize128,
			entropyInput:          "6bc709aa4c975b0eccb922ce2110fa9b572403f9013dfd10f06a88d54d380002",
			entropyInputReseed:    "cf1af84eddd5bef666ea42bea6067a23e52742e24661f944ba2514fe052abf31",
			additionalInputReseed: "a46988bad49b78c613c94e06a53b080bf6d20b7385bf4c782ad7cd145ddc9053",
			additionalInput1:      "139d6f72bf1d0ec5bfdd245e013f5cdac85e3eca716196018b92133c00a07436",
			additionalInput2:      "240f1a5af7fc2e4d32ace635acba5947f3564ecbfd7516c479c0adda20747f26",
			returnedBits: "e727268a546c0c891cf53a70a92820ee9bbf728ad52f30625b2e28f0f6c906f6" +
				"0ffd02f7d81623295950c04b63a48634eb41a5b4d649bdabff335ac3200690b0",
		},
		{
			name:               "AES-128 no df, PersonalizationStringLen = 256, AdditionalInputLen = 0",
			keySize:            KeySize128,
			entropyInput:       "34cbc2b217f3d907fa2ad6a0d7a813b0fda1e17fbeed94b0e0a0abfbec947146",
			personalization:    "e8fa4c5de825791e68180f2ba107e829c48299cb01be939cd0be76da120a91f2",
			entropyInputReseed: "8326f8e9cfbd02eb076bbb9819d96a02386f80bf913c8e4a80361d82cafad52e",
			returnedBits: "52f5e718bf48d99e498775c00378e545799bb2059aef0b74be573d8283f02b52" +
				"93917913bc8f26fc23760a1c86c3f5c844857419868eafeb17c9248227d026b8",
		},
		{
			name:                  "AES-192 no df, PersonalizationStringLen = 0, AdditionalInputLen = 320",
			keySize:               KeySize192,
			entropyInput:          "658a2df4d2e371a9a97192869f77ca77abc80e3f6d2e1e295d4d228de2ac90d54d6cbefa2e9ba3da",
			entropyInputReseed:    "79fba3e14d225b26e6365bd1ae1649ae93d716dc09ed509bbe338e2c78883c77aa0f43aed3c2bb2d",
			additionalInputReseed: "def5512314c47a6bb7cacb3529a66164ea48826650a54a6d9be9f5ee4fbb8c7d742d3deb005e8996",
			additionalInput1:      "54e3ad6f72c1f740a019572b7e446ad158e0a29bb0efbea60ff6f16596afbbcecfd35e65635f9f77",
			additionalInput2:      "6516e17a295d2fc141d2e8706e73822efd7b597fc576a97f43816cd1a9df4b9ff691324e2b581777",
			returnedBits: "932e59d41455f89ac8903019a683d5356894f4ad608f046f635779d3ae12e552" +
				"aa75575c7bbfdeb59d253893b96f71538fec6047582eb9388379ad939c85c956",
		},
		{
			name:               "AES-192 no df, PersonalizationStringLen = 320, AdditionalInputLen = 0",
			keySize:            KeySize192,
			entropyInput:       "e0e03d4a9e4de362e6bb9ba635b847912a3f9e1d83b4363fc258cdea2928a87823186f1c47c4175f",
			personalization:    "2eca97478d606b0bca56f055a1e8394b44f5e3028da59a3900d5d7bd341b0621b5f6d2309e36b4ae",
			entropyInputReseed: "66aa3db4bea8f2e1b2469751578f24f1b4d0df971cf2281712149812c20e804b90f3aac12d7c9832",
			returnedBits: "0e02f344b1159ccd7795974e7bb9df9e8cc381a130b98b43462285fe3ed85878" +
				"11abdef471bbcfe830e516b78f781a05b52e44d54122775084170ff6d80627e9",
		},
		{
			name:                  "AES-256 no df, PersonalizationStringLen = 0, AdditionalInputLen = 384",
			keySize:               KeySize256,
			entropyInput:          "99903165903fea49c2db26ed675e44cc14cb2c1f28b836b203240b02771e831146ffc4335373bb344688c5c950670291",
			entropyInputReseed:    "b4ee99fa9e0eddaf4a3612013cd636c4af69177b43eebb3c58a305b9979b68b5cc820504f6c029aad78a5d29c66e84a0",
			additionalInputReseed: "2d8c5c28b05696e74774eb69a10f01c5fabc62691ddf7848a8004bb5eeb4d2c5febe1aa01f4d557b23d7e9a0e4e90655",
			additionalInput1:      "0dc9cde42ac6e856f01a55f219c614de90c659260948db5053d414bab0ec2e13e995120c3eb5aafc25dc4bdcef8ace24",
			additionalInput2:      "711be6c035013189f362211889248ca8a3268e63a7eb26836d915810a680ac4a33cd1180811a31a0f44f08db3dd64f91",
			returnedBits: "11c7a0326ea737baa7a993d510fafee5374e7bbe17ef0e3e29f50fa68aac2124" +
				"b017d449768491cac06d136d691a4e80785739f9aaedf311bba752a3268cc531",
		},
		{
			name:               "AES-256 no df, PersonalizationStringLen = 384, AdditionalInputLen = 0",
			keySize:            KeySize256,
			entropyInput:       "ffad10100025a879672ff50374b286712f457dd01441d76ac1a1cd15c7390dd93179a2f5920d198bf34a1b76fbc21289",
			personalization:    "1d2be6f25e88fa30c4ef42e4d54efd957dec231fa00143ca47580be666a8c143a916c90b3819a0a7ea914e3c9a2e7a3f",
			entropyInputReseed: "6c1a089cae313363bc76a780139eb4f2f2048b1f6b07896c5c412bff0385440fc43b73facbb79e3a252fa01fe17ab391",
			returnedBits: "e053c7d4bd9099ef6a99f190a5fd80219437d642006672338da6e0fe73ca4d24" +
				"ffa51151bfbdac78d8a2f6255046edf57a04626e9977139c6933274299f3bdff",
		},
	}

	decode := func(s string) []byte {
		b, err := hex.DecodeString(s)
		is.NoError(err)
		return b
	}

	for _, v := range vectors {
		want := decode(v.returnedBits)

		cfg := DefaultConfig()
		cfg.KeySize = v.keySize
		cfg.Personalization = decode(v.personalization)
		cfg.BacktrackingResistance = true

		rand.Reader = bytes.NewReader(append(decode(v.entropyInput), decode(v.entropyInputReseed)...))
		d, err := newDRBG(&cfg)
		is.NoError(err, v.name)
		is.NoError(d.Reseed(decode(v.additionalInputReseed)), v.name)
		rand.Reader = originalReader

		got := make([]byte, len(want))
		_, err = d.ReadWithAdditionalInput(got, decode(v.additionalInput1))
		is.NoError(err, v.name)
		_, err = d.ReadWithAdditionalInput(got, decode(v.additionalInput2))
		is.NoError(err, v.name)
		is.Equal(want, got, v.name)
	}
}

// Test_DRBG_BacktrackingResistance_ReseedPaths verifies, with BacktrackingResistance enabled, that every automatic
// reseed (request count, reseed interval, fork detection, and key rotation) mixes fresh entropy into the existing key
// and counter exactly as an explicit NIST SP 800-90A §10.2.1.4.1 reseed does. For each path, two instances are
// created from identical entropy: one triggers the automatic reseed, the other performs the reseed directly, and
// their resulting key and counters must match. It also verifies that Reseed rejects additional input longer than key size + 16
// bytes.
//
// crypto/rand.Reader is temporarily replaced with known entropy, so this test must not run in parallel with other
// tests; the original reader is restored on exit.
func Test_DRBG_BacktrackingResistance_ReseedPaths(t *testing.T) {
	is := assert.New(t)

	originalReader := rand.Reader
	defer func() { rand.Reader = originalReader }()

	const seedLen = int(KeySize256) + 16
	entropy := append(bytes.Repeat([]byte{0x11}, seedLen), bytes.Repeat([]byte{0x22}, seedLen)...)

	paths := []struct {
		name    string
		mutate  func(cfg *Config)
		trigger func(d *drbg) error
	}{
		{
			name:   "request count",
			mutate: func(cfg *Config) { cfg.ReseedRequests = 1 },
			trigger: func(d *drbg) error {
				_, err := d.Read(make([]byte, 32))
				return err
			},
		},
		{
			name:   "reseed interval",
			mutate: func(cfg *Config) { cfg.ReseedInterval = time.Hour },
			trigger: func(d *drbg) error {
				d.lastReseedTime = time.Time{}
				_, err := d.Read(make([]byte, 32))
				return err
			},
		},
		{
			name:   "fork detection",
			mutate: func(cfg *Config) {},
			trigger: func(d *drbg) error {
				d.pid = -1
				_, err := d.Read(make([]byte, 32))
				return err
			},
		},
		{
			name:   "key rotation",
			mutate: func(cfg *Config) { cfg.EnableKeyRotation = true },
			trigger: func(d *drbg) error {
				d.asyncRekey()
				_, err := d.Read(make([]byte, 32))
				return err
			},
		},
	}

	for _, p := range paths {
		cfg := DefaultConfig()
		cfg.BacktrackingResistance = true
		p.mutate(&cfg)

		// The automatic path: one read, then the read that triggers the reseed.
		rand.Reader = bytes.NewReader(entropy)
		auto, err := newDRBG(&cfg)
		is.NoError(err, p.name)
		_, err = auto.Read(make([]byte, 32))
		is.NoError(err, p.name)
		is.NoError(p.trigger(auto), p.name)
		rand.Reader = originalReader

		// The reference: the same reads with an explicit conforming reseed in between.
		refCfg := DefaultConfig()
		refCfg.BacktrackingResistance = true
		rand.Reader = bytes.NewReader(entropy)
		ref, err := newDRBG(&refCfg)
		is.NoError(err, p.name)
		_, err = ref.Read(make([]byte, 32))
		is.NoError(err, p.name)
		is.NoError(ref.reseedUpdate(nil), p.name)
		_, err = ref.Read(make([]byte, 32))
		is.NoError(err, p.name)
		rand.Reader = originalReader

		// Compare immediately after the reseed: further reads could trigger additional reseeds.
		autoState, refState := auto.state.Load(), ref.state.Load()
		is.Equal(refState.key, autoState.key, "%s reseed must produce the same key as an explicit conforming reseed", p.name)
		is.Equal(refState.v, autoState.v, "%s reseed must produce the same counter as an explicit conforming reseed", p.name)
		is.Equal(ref.v, auto.v, "%s reseed must produce the same working counter as an explicit conforming reseed", p.name)
	}

	cfg := DefaultConfig()
	cfg.BacktrackingResistance = true
	d, err := newDRBG(&cfg)
	is.NoError(err)
	is.NoError(d.Reseed(make([]byte, seedLen)))
	is.Equal(ErrAdditionalInputTooLong, d.Reseed(make([]byte, seedLen+1)))
}

// Test_NewReader_BacktrackingResistance_ReseedRequests verifies that, with BacktrackingResistance enabled, NewReader
// enforces the NIST SP 800-90A §10.2.1 Table 3 maximum of 2^48 output requests between reseeds when ReseedRequests
// is not set, keeps an explicitly configured value, and leaves the default mode unchanged. It also verifies that a
// pooled instance reseeds once its request count reaches the maximum.
func Test_NewReader_BacktrackingResistance_ReseedRequests(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	r, err := NewReader(WithBacktrackingResistance(true))
	is.NoError(err)
	is.Equal(uint64(1<<48), r.Config().ReseedRequests, "unset ReseedRequests must default to 2^48")

	explicit, err := NewReader(WithBacktrackingResistance(true), WithReseedRequests(1000))
	is.NoError(err)
	is.Equal(uint64(1000), explicit.Config().ReseedRequests, "an explicit ReseedRequests must be kept")

	def, err := NewReader()
	is.NoError(err)
	is.Equal(uint64(0), def.Config().ReseedRequests, "default mode must leave ReseedRequests disabled")

	// A pooled instance at the maximum request count must reseed before its next output.
	rdr, ok := r.(*reader)
	is.True(ok)
	d, ok := rdr.pools[0].Get().(*drbg)
	is.True(ok)
	defer rdr.pools[0].Put(d)

	before := d.state.Load()
	atomic.StoreUint64(&d.requests, 1<<48)
	_, err = d.Read(make([]byte, 32))
	is.NoError(err)
	is.Equal(uint64(1), atomic.LoadUint64(&d.requests), "request count must restart after the reseed")
	is.NotSame(before, d.state.Load())
}

// scriptedBlock is a test cipher.Block whose Encrypt ignores its input and returns a scripted sequence of
// 16-byte blocks, used to simulate a stuck generator for the continuous health test.
type scriptedBlock struct {
	blocks [][16]byte
	next   *int
}

func (s scriptedBlock) BlockSize() int { return 16 }

func (s scriptedBlock) Encrypt(dst, _ []byte) {
	copy(dst, s.blocks[*s.next%len(s.blocks)][:])
	*s.next++
}

func (s scriptedBlock) Decrypt(dst, src []byte) { copy(dst, src) }

// Test_DRBG_ContinuousHealthTest verifies the FIPS 140-2 §4.9.2 continuous random number generator test: every
// generated block is compared with the previous one (within a request, across requests, and including the full
// block behind a partial final block), a failure clears the output and returns ErrHealthTestFailed, and the
// instance then stays in an error state, even after a reseed.
func Test_DRBG_ContinuousHealthTest(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	block := func(b byte) [16]byte { return [16]byte{0: b, 15: b} }

	cases := []struct {
		name     string
		script   [][16]byte
		requests []int
		failAt   int // index of the request expected to fail
	}{
		{"repeat within a request, not the first block", [][16]byte{block(1), block(2), block(3), block(3)}, []int{64}, 0},
		{"repeat across requests", [][16]byte{block(1), block(1)}, []int{16, 16}, 1},
		{"repeat in the partial final block", [][16]byte{block(1), block(2), block(2)}, []int{16, 17}, 1},
	}

	for _, c := range cases {
		cfg := DefaultConfig()
		cfg.ContinuousHealthTest = true
		d, err := newDRBG(&cfg)
		is.NoError(err, c.name)

		// Replace the cipher with the scripted one; the saved comparison block is from real output.
		next := 0
		d.state.Store(&state{block: scriptedBlock{blocks: c.script, next: &next}})

		for i, n := range c.requests {
			buf := bytes.Repeat([]byte{0xff}, n)
			got, err := d.Read(buf)
			if i < c.failAt {
				is.NoError(err, c.name)
				continue
			}
			is.Equal(ErrHealthTestFailed, err, c.name)
			is.Equal(0, got, c.name)
			is.Equal(make([]byte, n), buf, "%s: output must be cleared on failure", c.name)
		}

		// The error state persists, even after a reseed restores a working generator.
		is.NoError(d.Reseed(nil), c.name)
		_, err = d.Read(make([]byte, 32))
		is.Equal(ErrHealthTestFailed, err, "%s: error state must persist", c.name)
	}
}

// Test_DRBG_ContinuousHealthTest_FirstBlockWithheld verifies that, with the continuous health test enabled, the first
// block generated after an instance is created is not output (FIPS 140-2 §4.9.2): its first output equals the second
// block of an identically seeded instance without the test. It also verifies the test does not allocate.
//
// crypto/rand.Reader is temporarily replaced with known entropy, so this test must not run in parallel with other
// tests; the original reader is restored on exit.
func Test_DRBG_ContinuousHealthTest_FirstBlockWithheld(t *testing.T) {
	is := assert.New(t)

	originalReader := rand.Reader
	defer func() { rand.Reader = originalReader }()

	entropy := bytes.Repeat([]byte{0x3c}, int(KeySize256)+16)

	plainCfg := DefaultConfig()
	rand.Reader = bytes.NewReader(entropy)
	plain, err := newDRBG(&plainCfg)
	is.NoError(err)

	testedCfg := DefaultConfig()
	testedCfg.ContinuousHealthTest = true
	rand.Reader = bytes.NewReader(entropy)
	tested, err := newDRBG(&testedCfg)
	is.NoError(err)
	rand.Reader = originalReader

	want := make([]byte, 32)
	_, err = plain.Read(want)
	is.NoError(err)

	got := make([]byte, 16)
	_, err = tested.Read(got)
	is.NoError(err)
	is.Equal(want[16:], got, "the first generated block must be withheld")

	buf := make([]byte, 32)
	allocs := testing.AllocsPerRun(1000, func() { _, _ = tested.Read(buf) })
	is.Equal(0.0, allocs, "the continuous health test must not allocate")
}
