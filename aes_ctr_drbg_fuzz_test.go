// Copyright (c) 2024-2025 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package ctrdrbg

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Fuzz_Reader_Read exercises the Reader.Read method using a variety of buffer sizes.
// It ensures that Read does not return an error and produces the requested number of bytes
// for all valid sizes in the range [0, 65536]. Invalid sizes outside this range are skipped.
func Fuzz_Reader_Read(f *testing.F) {
	for _, sz := range []int{0, 1, 15, 16, 17, 32, 64, 1024, 4096} {
		f.Add(sz)
	}

	f.Fuzz(func(t *testing.T, size int) {
		is := assert.New(t)

		if size < 0 || size > 65536 {
			return // don't allocate insane slices
		}

		buf := make([]byte, size)
		n, err := Reader.Read(buf)
		is.NoError(err, "Reader.Read failed")
		is.Equal(size, n, "unexpected number of bytes read")
	})
}

// Fuzz_Reader_Concurrent tests the thread-safety of Reader.Read under concurrent access.
// For several buffer sizes, it spawns multiple goroutines that each perform a Read operation,
// checking that no errors occur. Sizes outside the range [1, 16384] are skipped.
func Fuzz_Reader_Concurrent(f *testing.F) {
	// Seed with a few realistic sizes
	for _, s := range []int{16, 1024, 4096, 16384} {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, size int) {
		is := assert.New(t)

		if size < 1 || size > 16384 {
			t.Skipf("slice size must be between 1 and 16384")
			return
		}

		const N = 8
		bufs := make([][]byte, N)
		errs := make(chan error, N)
		for i := range bufs {
			bufs[i] = make([]byte, size)
			go func(i int) {
				_, err := Reader.Read(bufs[i])
				errs <- err
			}(i)
		}
		for i := 0; i < N; i++ {
			err := <-errs
			is.NoError(err, "Concurrent Read failed")
		}
	})
}

// Fuzz_NewReader_AllOptions exercises NewReader with a variety of option combinations and parameter values.
// It fuzzes all tunable configuration fields, including key size, personalization, sharding, buffer size, and retry settings.
// If the key size is invalid, it asserts an error is returned. For valid configs, it checks that Read succeeds.
func Fuzz_NewReader_AllOptions(f *testing.F) {
	f.Add(uint64(32), int(3), int(5), int(0), int(1), int(16), true, []byte("seed"), int(16))
	f.Add(uint64(0), int(0), int(0), int(0), int(0), int(0), false, []byte(nil), int(1))
	f.Add(uint64(4096), int(10), int(10), int(5), int(32), int(0), true, []byte("p"), int(32))

	f.Fuzz(func(t *testing.T,
		maxBytes uint64,
		maxInitRetries int,
		maxRekeyAttempts int,
		shards int,
		bufSize int,
		keySizeRaw int,
		zeroBuffer bool,
		personalization []byte,
		mode int,
	) {
		is := assert.New(t)

		// Defensive bounds for fuzz
		if maxBytes > 1<<32 {
			maxBytes = 1 << 32
		}
		if bufSize < 0 {
			bufSize = 0
		}
		if bufSize > 1<<24 {
			bufSize = 1 << 24
		}
		if shards < 0 {
			shards = 0
		}
		if shards > 64 {
			shards = 64
		}
		if maxInitRetries < 0 {
			maxInitRetries = 0
		}
		if maxInitRetries > 100 {
			maxInitRetries = 100
		}
		if maxRekeyAttempts < 0 {
			maxRekeyAttempts = 0
		}
		if maxRekeyAttempts > 100 {
			maxRekeyAttempts = 100
		}
		if len(personalization) > 128 {
			personalization = personalization[:128]
		}

		// Choose a valid or invalid key size
		var keySize KeySize
		switch mode % 4 {
		case 0:
			keySize = KeySize128
		case 1:
			keySize = KeySize192
		case 2:
			keySize = KeySize256
		default:
			keySize = KeySize(keySizeRaw)
		}

		opts := []Option{
			WithMaxBytesPerKey(maxBytes),
			WithMaxInitRetries(maxInitRetries),
			WithMaxRekeyAttempts(maxRekeyAttempts),
			WithShards(shards),
			WithDefaultBufferSize(bufSize),
			WithEnableKeyRotation(mode%2 == 0),
			WithKeySize(keySize),
			WithUseZeroBuffer(zeroBuffer),
			WithPersonalization(personalization),
			WithRekeyBackoff(time.Duration(rand.Intn(1000)) * time.Millisecond),
			WithMaxRekeyBackoff(time.Duration(rand.Intn(3000)) * time.Millisecond),
		}

		r, err := NewReader(opts...)

		if keySize != KeySize128 && keySize != KeySize192 && keySize != KeySize256 {
			is.Error(err, "expected error with invalid keysize")
			return
		}

		// If we failed for any other reason (e.g. invalid config, entropy exhausted), just return.
		if err != nil {
			return
		}

		buf := make([]byte, 32)
		n, err := r.Read(buf)
		is.NoError(err, "Read failed")
		is.Equal(32, n, "short read")
	})
}

// Fuzz_NewReader_Personalization fuzzes NewReader with different personalization values to ensure
// that the DRBG accepts arbitrary domain separation strings. If NewReader returns an error (e.g. invalid config),
// the input is skipped. Otherwise, it asserts that a 16-byte read succeeds.
func Fuzz_NewReader_Personalization(f *testing.F) {
	f.Add([]byte("p"))
	f.Add([]byte{})
	f.Add([]byte(nil))
	f.Add(make([]byte, 64))

	f.Fuzz(func(t *testing.T, p []byte) {
		is := assert.New(t)

		r, err := NewReader(WithPersonalization(p))
		if err != nil {
			return // or t.Skip() to not count as "failure"
		}
		buf := make([]byte, 16)
		n, err := r.Read(buf)
		is.NoError(err, "read failed")
		is.Equal(16, n, "short read")
	})
}

// Fuzz_NewReader_Buffers fuzzes NewReader with various buffer size configurations, verifying
// that the reader is correctly initialized and produces the expected output. If initialization fails
// (due to an unsupported buffer size, etc.), the input is skipped. Otherwise, it checks a 32-byte read.
func Fuzz_NewReader_Buffers(f *testing.F) {
	for _, sz := range []int{0, 1, 16, 1024, 1 << 20, 1 << 23} {
		f.Add(sz)
	}
	f.Fuzz(func(t *testing.T, bufSize int) {
		is := assert.New(t)

		if bufSize < 0 || bufSize > 1<<24 {
			return
		}
		r, err := NewReader(WithDefaultBufferSize(bufSize))
		if err != nil {
			return // or t.Skip() to not count as "failure"
		}
		buf := make([]byte, 32)
		n, err := r.Read(buf)
		is.NoError(err, "read failed")
		is.Equal(32, n, "short read")
	})
}

func Fuzz_ReadWithAdditionalInput(f *testing.F) {
	f.Add(32, []byte("entropy1"))
	f.Add(16, []byte(nil))
	f.Add(64, []byte("another-entropy-value"))

	f.Fuzz(func(t *testing.T, bufSize int, addIn []byte) {
		is := assert.New(t)

		if bufSize < 0 || bufSize > 4096 {
			return
		}
		r, err := NewReader()
		if err != nil {
			return
		}
		buf := make([]byte, bufSize)
		n, err := r.ReadWithAdditionalInput(buf, addIn)
		is.NoError(err, "ReadWithAdditionalInput failed")
		is.Equal(bufSize, n, "short read")
		if bufSize > 0 {
			// Ensure not all zeros (statistically, not guaranteed but good quick check)
			zeros := true
			for _, b := range buf {
				if b != 0 {
					zeros = false
					break
				}
			}
			is.False(zeros, "all zero output")
		}
	})
}

func Fuzz_Reseed_Concurrency(f *testing.F) {
	f.Add([]byte("input1"))
	f.Add([]byte{})
	f.Fuzz(func(t *testing.T, addIn []byte) {
		r, err := NewReader()
		if err != nil {
			return
		}
		const n = 10
		done := make(chan struct{}, n)
		for i := 0; i < n; i++ {
			go func() {
				_ = r.Reseed(addIn)
				done <- struct{}{}
			}()
		}
		for i := 0; i < n; i++ {
			<-done
		}
		// No assertion: not crashing is the property.
	})
}

func Fuzz_Counter_Overflow(f *testing.F) {
	// Add edge cases: (idx, val)
	f.Add(15, byte(255))
	f.Add(0, byte(0))
	f.Add(7, byte(127))
	f.Add(3, byte(1))

	f.Fuzz(func(t *testing.T, idx int, val byte) {
		is := assert.New(t)
		cfg := DefaultConfig()
		d, err := newDRBG(&cfg)
		if err != nil {
			return
		}
		if idx >= 0 && idx < len(d.v) {
			// Set counter to all 0xff, then set one index to val
			for i := range d.v {
				d.v[i] = 0xff
			}
			d.v[idx] = val
			buf := make([]byte, 16)
			_, err := d.Read(buf)
			is.NoError(err)
		}
	})
}

func Fuzz_Config_FunctionalOptions_Combinatorics(f *testing.F) {
	f.Add(32, true, true, 128, 7, 42, 100, []byte("domain"), 1)
	f.Fuzz(func(t *testing.T,
		bufSz int, rot, zero bool, keySz int, maxInit, maxRekey, rekeyB int, pers []byte, mode int,
	) {
		is := assert.New(t)

		if bufSz < 0 || bufSz > 1<<20 {
			return
		}
		var k KeySize
		switch keySz % 3 {
		case 0:
			k = KeySize128
		case 1:
			k = KeySize192
		default:
			k = KeySize256
		}
		opts := []Option{
			WithKeySize(k),
			WithEnableKeyRotation(rot),
			WithUseZeroBuffer(zero),
			WithDefaultBufferSize(bufSz),
			WithMaxInitRetries(maxInit),
			WithMaxRekeyAttempts(maxRekey),
			WithRekeyBackoff(time.Duration(rekeyB) * time.Millisecond),
			WithMaxRekeyBackoff(time.Duration(maxRekey) * time.Millisecond),
			WithPersonalization(pers),
		}
		r, err := NewReader(opts...)
		if err != nil {
			return
		}
		buf := make([]byte, 64)
		n, err := r.Read(buf)
		is.NoError(err)
		is.Equal(64, n)
	})
}
