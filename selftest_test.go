// Copyright (c) 2024-2026 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package ctrdrbg

import (
	"strings"
	"sync"
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

// Test_RunDRBGKAT_Success verifies that the NIST SP 800-90A §11.3 known-answer tests of the instantiate, reseed,
// and generate functions pass. The tests are run directly because the result used by NewReader is cached per process.
func Test_RunDRBGKAT_Success(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	is.NoError(runDRBGKAT(), "DRBG known-answer tests should pass with NIST CAVP vectors")
}

// Test_NewReader_BacktrackingResistance_RunsDRBGSelfTests verifies that NewReader runs the DRBG known-answer tests
// when BacktrackingResistance is enabled and records a passing result.
func Test_NewReader_BacktrackingResistance_RunsDRBGSelfTests(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	_, err := NewReader(WithBacktrackingResistance(true))
	is.NoError(err)
	is.NoError(runDRBGSelfTests(), "the cached DRBG self-test result should be a pass")
}

// Test_NewReader_BacktrackingResistance_DRBGSelfTestFailure verifies that NewReader with BacktrackingResistance fails
// with ErrDRBGSelfTestFailed when a DRBG known-answer test fails, and that the failure persists for later calls.
//
// It replaces a known answer and resets the per-process self-test result, so it must not run in parallel with other
// tests; the original vectors and a fresh result are restored on exit.
func Test_NewReader_BacktrackingResistance_DRBGSelfTestFailure(t *testing.T) {
	is := assert.New(t)

	original := drbgKATVectors
	defer func() {
		drbgKATVectors = original
		drbgSelfTestOnce = sync.Once{}
		drbgSelfTestErr = nil
	}()

	corrupted := append([]drbgKATVector(nil), original...)
	corrupted[0].returnedBits = strings.Repeat("00", 64)
	drbgKATVectors = corrupted
	drbgSelfTestOnce = sync.Once{}
	drbgSelfTestErr = nil

	_, err := NewReader(WithBacktrackingResistance(true))
	is.Equal(ErrDRBGSelfTestFailed, err, "NewReader must fail when a DRBG known-answer test fails")

	// The failure is kept for the life of the process, even if the vectors are restored.
	drbgKATVectors = original
	_, err = NewReader(WithBacktrackingResistance(true))
	is.Equal(ErrDRBGSelfTestFailed, err, "the self-test failure must persist")

	_, err = NewReader()
	is.NoError(err, "the default mode does not run the DRBG known-answer tests")
}
