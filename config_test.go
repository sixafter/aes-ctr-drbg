// Copyright (c) 2024-2026 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package ctrdrbg

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestConfig_DefaultConfig checks that DefaultConfig returns a Config with documented default values.
func TestConfig_DefaultConfig(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	cfg := DefaultConfig()
	is.Equal(KeySize256, cfg.KeySize, "KeySize should default to 32 (AES-256)")
	is.Equal(uint64(1<<30), cfg.MaxBytesPerKey, "MaxBytesPerKey should default to 1GiB")
	is.Equal(3, cfg.MaxInitRetries, "MaxInitRetries should default to 3")
	is.Equal(5, cfg.MaxRekeyAttempts, "MaxRekeyAttempts should default to 5")
	is.Equal(2*time.Second, cfg.MaxRekeyBackoff, "MaxRekeyBackoff should default to 2s")
	is.Equal(100*time.Millisecond, cfg.RekeyBackoff, "RekeyBackoff should default to 100ms")
	is.False(cfg.EnableKeyRotation, "EnableKeyRotation should default to false")
	is.Nil(cfg.Personalization, "Personalization should default to nil")
}

// TestConfig_WithKeySize verifies that WithKeySize correctly overrides the KeySize field.
func TestConfig_WithKeySize(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	cfg := DefaultConfig()
	WithKeySize(KeySize128)(&cfg)
	is.Equal(KeySize128, cfg.KeySize, "WithKeySize should override KeySize")
}

// TestConfig_WithMaxBytesPerKey verifies that WithMaxBytesPerKey correctly overrides the MaxBytesPerKey field.
func TestConfig_WithMaxBytesPerKey(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	cfg := DefaultConfig()
	WithMaxBytesPerKey(42)(&cfg)
	is.Equal(uint64(42), cfg.MaxBytesPerKey, "WithMaxBytesPerKey should override MaxBytesPerKey")
}

// TestConfig_WithMaxInitRetries ensures that WithMaxInitRetries updates only the MaxInitRetries field.
func TestConfig_WithMaxInitRetries(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	cfg := DefaultConfig()
	WithMaxInitRetries(7)(&cfg)
	is.Equal(7, cfg.MaxInitRetries, "WithMaxInitRetries should override MaxInitRetries")
}

// TestConfig_WithMaxRekeyAttempts checks that WithMaxRekeyAttempts sets MaxRekeyAttempts correctly.
func TestConfig_WithMaxRekeyAttempts(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	cfg := DefaultConfig()
	WithMaxRekeyAttempts(11)(&cfg)
	is.Equal(11, cfg.MaxRekeyAttempts, "WithMaxRekeyAttempts should override MaxRekeyAttempts")
}

// TestConfig_WithMaxRekeyBackoff verifies that WithMaxRekeyBackoff updates the MaxRekeyBackoff field.
func TestConfig_WithMaxRekeyBackoff(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	cfg := DefaultConfig()
	WithMaxRekeyBackoff(888 * time.Millisecond)(&cfg)
	is.Equal(888*time.Millisecond, cfg.MaxRekeyBackoff, "WithMaxRekeyBackoff should override MaxRekeyBackoff")
}

// TestConfig_WithRekeyBackoff checks that WithRekeyBackoff sets the RekeyBackoff field.
func TestConfig_WithRekeyBackoff(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	cfg := DefaultConfig()
	WithRekeyBackoff(222 * time.Millisecond)(&cfg)
	is.Equal(222*time.Millisecond, cfg.RekeyBackoff, "WithRekeyBackoff should override RekeyBackoff")
}

// TestConfig_WithEnableKeyRotation validates WithEnableKeyRotation toggles EnableKeyRotation as expected.
func TestConfig_WithEnableKeyRotation(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	cfg := DefaultConfig()
	WithEnableKeyRotation(false)(&cfg)
	is.False(cfg.EnableKeyRotation, "WithEnableKeyRotation(false) should set EnableKeyRotation to false")
	WithEnableKeyRotation(true)(&cfg)
	is.True(cfg.EnableKeyRotation, "WithEnableKeyRotation(true) should set EnableKeyRotation to true")
}

// TestConfig_WithPersonalization checks that WithPersonalization sets the Personalization field.
func TestConfig_WithPersonalization(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	cfg := DefaultConfig()
	pers := []byte("unique-domain")
	WithPersonalization(pers)(&cfg)
	is.Equal(pers, cfg.Personalization, "WithPersonalization should set Personalization")
}

// TestConfig_WithUseZeroBuffer checks that WithUseZeroBuffer sets the UseZeroBuffer field correctly.
func TestConfig_WithUseZeroBuffer(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	cfg := DefaultConfig()
	WithUseZeroBuffer(true)(&cfg)
	is.True(cfg.UseZeroBuffer, "WithUseZeroBuffer(true) should set UseZeroBuffer to true")
	WithUseZeroBuffer(false)(&cfg)
	is.False(cfg.UseZeroBuffer, "WithUseZeroBuffer(false) should set UseZeroBuffer to false")
}

// TestConfig_WithDefaultBufferSize verifies that WithDefaultBufferSize sets the DefaultBufferSize field correctly.
func TestConfig_WithDefaultBufferSize(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	cfg := DefaultConfig()
	WithDefaultBufferSize(64)(&cfg)
	is.Equal(64, cfg.DefaultBufferSize, "WithDefaultBufferSize should set DefaultBufferSize")
}

// TestConfig_WithShards ensures that WithShards updates only the Shards field.
func TestConfig_WithShards(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	cfg := DefaultConfig()
	WithShards(8)(&cfg)
	is.Equal(8, cfg.Shards, "WithShards should override Shards")
}

// TestConfig_CombinedOptions ensures multiple options can be applied sequentially.
func TestConfig_CombinedOptions(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	cfg := DefaultConfig()
	opts := []Option{
		WithKeySize(KeySize192),
		WithMaxBytesPerKey(1024),
		WithMaxInitRetries(2),
		WithMaxRekeyAttempts(8),
		WithMaxRekeyBackoff(345 * time.Millisecond),
		WithRekeyBackoff(123 * time.Millisecond),
		WithEnableKeyRotation(false),
		WithPersonalization([]byte("tenant42")),
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	is.Equal(KeySize192, cfg.KeySize)
	is.Equal(uint64(1024), cfg.MaxBytesPerKey)
	is.Equal(2, cfg.MaxInitRetries)
	is.Equal(8, cfg.MaxRekeyAttempts)
	is.Equal(345*time.Millisecond, cfg.MaxRekeyBackoff)
	is.Equal(123*time.Millisecond, cfg.RekeyBackoff)
	is.False(cfg.EnableKeyRotation)
	is.Equal([]byte("tenant42"), cfg.Personalization)
}

// TestConfig_WithReseedInterval verifies that WithReseedInterval sets the ReseedInterval field.
func TestConfig_WithReseedInterval(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	cfg := DefaultConfig()
	WithReseedInterval(5 * time.Second)(&cfg)
	is.Equal(5*time.Second, cfg.ReseedInterval, "WithReseedInterval should set ReseedInterval")
}

// TestConfig_WithReseedRequests verifies that WithReseedRequests sets the ReseedRequests field.
func TestConfig_WithReseedRequests(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	cfg := DefaultConfig()
	WithReseedRequests(42)(&cfg)
	is.Equal(uint64(42), cfg.ReseedRequests, "WithReseedRequests should set ReseedRequests")
}

// TestConfig_WithPredictionResistance verifies that WithPredictionResistance sets the PredictionResistance field.
func TestConfig_WithPredictionResistance(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	cfg := DefaultConfig()
	WithPredictionResistance(true)(&cfg)
	is.True(cfg.PredictionResistance, "WithPredictionResistance(true) should set PredictionResistance to true")
	WithPredictionResistance(false)(&cfg)
	is.False(cfg.PredictionResistance, "WithPredictionResistance(false) should set PredictionResistance to false")
}

func TestConfig_WithForkDetectionInterval(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	cfg := DefaultConfig()
	WithForkDetectionInterval(42)(&cfg)
	is.Equal(uint64(42), cfg.ForkDetectionInterval, "WithForkDetectionInterval should set ForkDetectionInterval")
}

// TestConfig_WithReseedRequests_Clamping verifies that WithReseedRequests clamps values exceeding 2^48.
func TestConfig_WithReseedRequests_Clamping(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	cfg := DefaultConfig()

	// Value within limit should be set as-is
	WithReseedRequests(1000)(&cfg)
	is.Equal(uint64(1000), cfg.ReseedRequests, "Values within limit should be set as-is")

	// Value at exactly 2^48 should be set as-is
	WithReseedRequests(1 << 48)(&cfg)
	is.Equal(uint64(1<<48), cfg.ReseedRequests, "Value at 2^48 should be set as-is")

	// Value exceeding 2^48 should be clamped
	WithReseedRequests((1 << 48) + 1)(&cfg)
	is.Equal(uint64(1<<48), cfg.ReseedRequests, "Value exceeding 2^48 should be clamped to 2^48")

	// Very large value should be clamped
	WithReseedRequests(^uint64(0))(&cfg) // max uint64
	is.Equal(uint64(1<<48), cfg.ReseedRequests, "Max uint64 should be clamped to 2^48")
}

// TestConfig_WithSelfTests verifies that WithSelfTests sets the EnableSelfTests field.
func TestConfig_WithSelfTests(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	cfg := DefaultConfig()
	is.False(cfg.EnableSelfTests, "EnableSelfTests should default to false")

	WithSelfTests(true)(&cfg)
	is.True(cfg.EnableSelfTests, "WithSelfTests(true) should set EnableSelfTests to true")

	WithSelfTests(false)(&cfg)
	is.False(cfg.EnableSelfTests, "WithSelfTests(false) should set EnableSelfTests to false")
}

// TestConfig_WithZeroization verifies that WithZeroization sets the EnableZeroization field.
func TestConfig_WithZeroization(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	cfg := DefaultConfig()
	is.False(cfg.EnableZeroization, "EnableZeroization should default to false")

	WithZeroization(true)(&cfg)
	is.True(cfg.EnableZeroization, "WithZeroization(true) should set EnableZeroization to true")

	WithZeroization(false)(&cfg)
	is.False(cfg.EnableZeroization, "WithZeroization(false) should set EnableZeroization to false")
}

// TestConfig_WithContinuousHealthTest verifies that WithContinuousHealthTest sets the ContinuousHealthTest field.
func TestConfig_WithContinuousHealthTest(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	cfg := DefaultConfig()
	is.False(cfg.ContinuousHealthTest, "ContinuousHealthTest should default to false")

	WithContinuousHealthTest(true)(&cfg)
	is.True(cfg.ContinuousHealthTest, "WithContinuousHealthTest(true) should set ContinuousHealthTest to true")

	WithContinuousHealthTest(false)(&cfg)
	is.False(cfg.ContinuousHealthTest, "WithContinuousHealthTest(false) should set ContinuousHealthTest to false")
}
