// Copyright (c) 2024-2025 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package ctrdrbg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_RunSelfTests_Success verifies that RunSelfTests passes with a valid AES-CTR implementation.
func Test_RunSelfTests_Success(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// RunSelfTests uses sync.Once, so we test the underlying runKAT directly
	// to ensure consistent test behavior across runs.
	err := runKAT()
	is.NoError(err, "runKAT should pass with valid NIST test vector")
}

// Test_NewReader_WithSelfTests_Enabled verifies that NewReader runs self-tests when enabled.
func Test_NewReader_WithSelfTests_Enabled(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	rdr, err := NewReader(WithSelfTests(true))
	is.NoError(err, "NewReader should succeed when self-tests pass")
	is.NotNil(rdr, "Reader should be returned")

	// Verify the config reflects the option
	cfg := rdr.Config()
	is.True(cfg.EnableSelfTests, "EnableSelfTests should be true")
}

// Test_NewReader_WithSelfTests_Disabled verifies that NewReader skips self-tests when disabled.
func Test_NewReader_WithSelfTests_Disabled(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	rdr, err := NewReader(WithSelfTests(false))
	is.NoError(err, "NewReader should succeed without self-tests")
	is.NotNil(rdr, "Reader should be returned")

	// Verify the config reflects the option
	cfg := rdr.Config()
	is.False(cfg.EnableSelfTests, "EnableSelfTests should be false")
}

// Test_NewReader_SelfTests_DefaultDisabled verifies that self-tests are disabled by default.
func Test_NewReader_SelfTests_DefaultDisabled(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	rdr, err := NewReader()
	is.NoError(err, "NewReader should succeed with defaults")

	cfg := rdr.Config()
	is.False(cfg.EnableSelfTests, "EnableSelfTests should default to false")
}
