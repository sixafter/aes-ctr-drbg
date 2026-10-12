// Copyright (c) 2024-2026 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

//go:build !windows

package ctrdrbg

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_DRBG_ForkDetectionInterval_WithReseedRequests ensures interval-based fork detection still
// runs when ReseedRequests is smaller than ForkDetectionInterval (reseeds must not reset the fork-check counter).
func Test_DRBG_ForkDetectionInterval_WithReseedRequests(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	cfg := DefaultConfig()
	cfg.ReseedRequests = 2
	cfg.ForkDetectionInterval = 4
	d, err := newDRBG(&cfg)
	is.NoError(err)

	// Simulate a fork by caching a PID that differs from the current process.
	d.pid = -1

	buf := make([]byte, 32)
	for i := 0; i < int(cfg.ForkDetectionInterval); i++ {
		_, err = d.Read(buf)
		is.NoError(err)
	}

	is.Equal(os.Getpid(), d.pid, "fork should be detected after ForkDetectionInterval requests")
}

// Test_DRBG_ForkReseedFailure_Retried ensures that when the post-fork reseed fails, the cached PID is not
// updated, so the reseed is retried at the next fork check instead of the child keeping the parent's state.
func Test_DRBG_ForkReseedFailure_Retried(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	cfg := DefaultConfig()
	d, err := newDRBG(&cfg)
	is.NoError(err)

	// Simulate a fork by caching a PID that differs from the current process.
	d.pid = -1

	// Force the post-fork reseed to fail: an invalid key size makes cipher construction fail.
	cfg.KeySize = 15
	buf := make([]byte, 32)
	_, err = d.Read(buf)
	is.NoError(err)
	is.Equal(-1, d.pid, "cached PID must not be updated when the post-fork reseed fails")

	// Once reseeding can succeed again, the next fork check must reseed and record the new PID.
	cfg.KeySize = KeySize256
	_, err = d.Read(buf)
	is.NoError(err)
	is.Equal(os.Getpid(), d.pid, "cached PID should be updated after a successful post-fork reseed")
}
