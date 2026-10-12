# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Date format: `YYYY-MM-DD`

---
## [Unreleased]

### Added
### Changed
### Deprecated
### Removed
### Fixed
### Security

---

## [1.21.0] - 2026-10-11

### Added
- **feature:** Added `WithBacktrackingResistance` option to enable the NIST SP 800-90A §10.2.1.5.1 (step 6) post-generate update, which derives a new key and counter after every output request for backtracking resistance; when enabled, instantiation also follows §10.2.1.3.1 (initial key and counter derived from the seed material by the §10.2.1.2 update function), verified against NIST CAVP CTR_DRBG known-answer vectors (opt-in; `2 allocs/op` when enabled, default output path unchanged at `0 allocs/op`).
- **feature:** Added `ErrPersonalizationTooLong`, returned by `NewReader` when `BacktrackingResistance` is enabled and the personalization string exceeds the NIST SP 800-90A §10.2.1 Table 3 maximum of key size + 16 bytes.
- **feature:** With `BacktrackingResistance` enabled, `ReadWithAdditionalInput` now mixes additional input into the key and counter before and after generation per NIST SP 800-90A §10.2.1.5.1 steps 2 and 6 (without obtaining new entropy), treats empty additional input as absent, and returns the new `ErrAdditionalInputTooLong` when additional input exceeds key size + 16 bytes; verified against NIST CAVP CTR_DRBG known-answer vectors (`ReadWithAdditionalInput` with additional input: `4 allocs/op`, previously `5 allocs/op`).
- **feature:** With both `PredictionResistance` and `BacktrackingResistance` enabled, the reseed before each output request now follows NIST SP 800-90A §9.3.1 step 7 and §10.2.1.4.1: fresh entropy is combined with the request's additional input and mixed into the existing key and counter, and the additional input is not used again for that request; verified against NIST CAVP CTR_DRBG prediction-resistance known-answer vectors (allocs/op unchanged in every mode).
- **feature:** With `BacktrackingResistance` enabled, every reseed (explicit `Reseed`, `ReseedInterval`, `ReseedRequests`, fork detection, and automatic key rotation) now follows NIST SP 800-90A §10.2.1.4.1: fresh entropy combined with the optional additional input is mixed into the existing key and counter, and the personalization string is not used; `Reseed` returns `ErrAdditionalInputTooLong` when additional input exceeds key size + 16 bytes; verified against NIST CAVP CTR_DRBG reseed known-answer vectors (allocs/op unchanged in every mode).
- **feature:** With `BacktrackingResistance` enabled, `NewReader` enforces the NIST SP 800-90A §10.2.1 Table 3 limit of 2^48 output requests between reseeds: an unset `ReseedRequests` defaults to 2^48, an explicit value is kept, and the default mode is unchanged (allocs/op unchanged in every mode).
- **feature:** With `BacktrackingResistance` enabled, `NewReader` runs NIST SP 800-90A §11.3 known-answer self-tests of the instantiate, reseed, and generate functions (with and without prediction resistance, for each key size) once per process before first use, using NIST CAVP CTR_DRBG vectors; a failure returns the new `ErrDRBGSelfTestFailed` and persists for the life of the process (read-path allocs/op unchanged).

### Changed
- **debt:** Upgraded dependencies to their latest stable versions.
- **debt:** Documented that, with `EnableKeyRotation` enabled, a concurrent asynchronous key rotation may replace state derived from `ReadWithAdditionalInput` additional input.
- **debt:** Modified [README.md](../README.md) to describe NIST SP 800-90A conformance as opt-in via `WithBacktrackingResistance`, to describe the conforming prediction-resistance reseed, and to rename the "NIST SP 800-90A Compliance" section to "NIST SP 800-90A Specification Mapping".
- **debt:** Rewrote the [README.md](../README.md) overview and feature descriptions to state verifiable behavior (allocations, options, concurrency, fork detection, testing) instead of general claims, and corrected the statement that output is deterministic for a given seed: output is determined by the entropy read from `crypto/rand.Reader`, the personalization string, and any additional input.
- **debt:** Modified [docs/NIST-SP-800-90A.md](../docs/NIST-SP-800-90A.md) to map each NIST SP 800-90A requirement to the default mode, the `WithBacktrackingResistance` mode, and the tests that verify it, with separate sections for related FIPS 140-2 features and other behavior.

### Deprecated
### Removed
### Fixed
- **defect:** Fixed documentation comment related to the default value for `EnableKeyRotation`.
- **defect:** Fixed a panic (`crypto/aes: output not full block`) when `UseZeroBuffer` is enabled and the requested length is not a multiple of the AES block size (16 bytes).
- **defect:** Fixed interval-based fork detection (`ForkDetectionInterval > 0`) never running when combined with `ReseedRequests`; fork checks now use a dedicated counter that reseeds do not reset.
- **defect:** Fixed fork detection permanently accepting a failed post-fork reseed; the cached PID is now updated only after a successful reseed, so a failed reseed is retried at the next fork check.
- **defect:** Fixed documentation comment for `Reseed`, which reseeds one DRBG instance per shard pool rather than all pooled instances.
- **defect:** Fixed documentation for `EnableZeroization`/`WithZeroization` and the NIST SP 800-90A mapping to describe zeroization as best-effort during automatic key rotation; the AES key schedule held by `cipher.Block`, seed buffers, and state replaced by reseeds are not erased.
- **defect:** Fixed `ContinuousHealthTest` to comply with the FIPS 140-2 §4.9.2 continuous random number generator test: every generated 128-bit block (including the full block behind a partial final block) is compared with the previous one on the same pooled instance across requests, the first block generated after an instance is created is withheld for comparison, and a failure clears the output buffer and leaves the instance in a permanent error state. Previously only the first block of each request was compared and requests shorter than 16 bytes were not tested (still opt-in; `0 allocs/op`).
- **defect:** Fixed the continuous health test citation, which referred to NIST SP 800-90A §11.3.3 ("Testing the Generate Function"); the test is the FIPS 140-2 §4.9.2 continuous random number generator test.
- **defect:** Fixed `DefaultConfig` documentation that described the defaults as compliant with FIPS 140-2 and NIST SP 800-90A; NIST SP 800-90A §10.2.1 conformance is opt-in via `WithBacktrackingResistance`.
- **defect:** Fixed the package documentation, which stated that the package follows the NIST SP 800-90A §10.2.1 construction in all modes; it now describes the default mode and the conforming mode enabled by `WithBacktrackingResistance`.
- **defect:** Fixed flaky `Test_CTRDRBG_ReadWithDifferentBufferSizes` and `Fuzz_ReadWithAdditionalInput`, which asserted that 1-byte random output is not all zeros (a 1 in 256 chance); the check now applies to buffers of at least 8 bytes, and subtest names show the buffer size.
- **defect:** Fixed newly constructed DRBG instances performing an unnecessary reseed on their first output request when `ReseedInterval` is set; the interval is now measured from instantiation.
- **defect:** Fixed documentation stating that prediction resistance ignores additional input "as required by the standard"; NIST SP 800-90A §9.3.1 mixes it into the prediction-resistance reseed. The default mode continues to ignore it, and the documentation now says so.
- **defect:** Fixed `Reseed` documentation describing the method as compliant with NIST SP 800-90A in all modes; it now describes the default reseed and the §10.2.1.4.1 reseed used with `BacktrackingResistance`.

### Security

---

## [1.20.0] - 2026-08-22

### Added
### Changed
- **debt:** Upgraded to [Go 1.27](https://go.dev/doc/go1.27).
- **debt:** Upgraded dependencies to their latest stable versions.

### Deprecated
### Removed
### Fixed
### Security

---

## [1.19.3] - 2026-08-13

### Added
### Changed
- **debt:** Upgraded dependencies to latest stable versions.

### Deprecated
### Removed
### Fixed
### Security

---

## [1.19.2] - 2026-06-29

### Added
### Changed
- **debt:** Upgraded dependencies to latest stable versions.

### Deprecated
### Removed
### Fixed
### Security

---

## [1.19.1] - 2026-05-08

### Added
### Changed
- **debt:** Upgraded dependencies to latest stable versions.

### Deprecated
### Removed
### Fixed
### Security

---

## [1.19.0] - 2026-03-26

### Added
### Changed
- **debt:** Upgraded [Cosign](https://github.com/sigstore/cosign-installer) to latest stable version.

### Deprecated
### Removed
### Fixed
### Security

---

## [1.18.0] - 2026-03-14


### Added
### Changed
- **debt:** Upgraded all dependencies to their latest stable versions.

### Deprecated
### Removed
### Fixed
### Security

---

## [1.17.0] - 2026-02-12

### Added
### Changed
- **debt:** Upgraded all dependencies to their latest stable versions, which includes Go 1.26.
### Deprecated
### Removed
### Fixed
### Security

---

## [1.16.0] - 2026-01-22


### Added
- **feature:** Added `WithContinuousHealthTest` option to enable NIST SP 800-90A §11.3.3 continuous health testing for FIPS 140-2/140-3 compliance (detects stuck DRBG output).

### Changed
- **debt:** Minor clarifications in [README.md](../README.md) regarding FIPS 140-2 compliance features.
- **debt:** Upgraded copyright notice to include 2026.

### Deprecated
### Removed
### Fixed
### Security

---

## [1.15.0] - 2025-12-29


### Added
- **feature:** Added `WithSelfTests` option to enable FIPS 140-2 Known Answer Test (KAT) self-tests at DRBG initialization using NIST CAVP vectors.
- **feature:** Added `WithZeroization` option to enable secure erasure of old key material during key rotation for forward secrecy (FIPS 140-2 §4.7.6 compliance).

### Changed
- **risk:** Added 64 KB maximum request size limit per DRBG call to comply with NIST SP 800-90A §10.2.1.
- **risk:** Clamped user-configured reseed interval to a maximum of 2^48 requests per NIST SP 800-90A §10.2.1.1.
- **debt:** Updated NIST SP 800-90A specification mapping table in documentation to reflect new FIPS compliance features (max request size, reseed interval cap, KATs, and key zeroization).

### Deprecated
### Removed
### Fixed
### Security

---

## [1.14.5] - 2025-11-21

### Added
- **debt:** Modify README with how AES-CTR-DRBG is used when FIPS mode is enabled.

### Changed
### Deprecated
### Removed
### Fixed
- **defect:** Fix issue [#24](https://github.com/sixafter/aes-ctr-drbg/issues/24) by adding verification to ensure the Go module ZIP served via `GOPROXY=direct` matches the ZIP served by `proxy.golang.org`.

### Security

---

## [1.14.1] - 2025-11-20

### Added
- **risk**: Added signature verification make target to match the README instructions.
- **risk:** Add go module verification make target to verify module checksums.

### Changed
### Deprecated
### Removed
### Fixed
- **defect:** Fixed `README.md` instructions for verifying module checksums.

### Security

---

## [1.14.0] - 2025-11-20

### Added
### Changed
- **debt:** Upgraded all dependencies to their latest stable versions.

### Deprecated
### Removed
### Fixed
- **defect:** Fix [Issue #24](https://github.com/sixafter/aes-ctr-drbg/issues/24): v1.13.0 checksum conflict in golang proxy breaks build.
### Security

---

## [1.13.0] - 2025-11-07

### Added
### Changed
- **debt:** Upgraded all dependencies to their latest stable versions.

### Deprecated
### Removed
### Fixed
### Security

---

## [1.12.0] - 2025-10-16

### Added
### Changed
- **debt:** Upgraded all dependencies to their latest stable versions.
- **debt:** Updated documentation and Go-doc comments.

### Deprecated
### Removed
### Fixed
### Security

---

## [1.11.0] - 2025-10-08

### Added
### Changed
- **debt:** Upgraded all dependencies to their latest stable versions.
- **debt:** Updated documentation and Go-doc comments.
- **debt:** Replaced panic-based failure handling in `initShardPools` (called by `NewReader()`) with explicit error propagation.

### Deprecated
### Removed
### Fixed
### Security

---

## [1.10.0] - 2025-09-30

### Added
### Changed
- **debt:** Upgraded all dependencies to their latest stable versions.

### Deprecated
### Removed
### Fixed
### Security

---

## [1.9.0] - 2025-09-15

### Added
### Changed
- **debt:** Upgraded all dependencies to their latest stable versions.

### Deprecated
### Removed
### Fixed
### Security

---

## [1.8.0] - 2025-09-10

### Added
### Changed
- **debt:** Upgraded all dependencies to their latest stable versions.

### Deprecated
### Removed
### Fixed
### Security

---

## [1.7.0] - 2025-09-01

### Added
### Changed
- **debt:** Upgraded all dependencies to their latest stable versions.
- **risk:** Updated copyright to reflect date range through present year.
- **risk:** Removed `t.Parallel()` from fuzz tests as fuzzing already runs inputs in parallel. 
  - The fuzz engine schedules many inputs concurrently across workers (bounded by CPUs / `-test.parallel`). Adding `t.Parallel()` makes each input’s subtest itself run in parallel with others—an unnecessary second layer of parallelism.

### Deprecated
### Removed
### Fixed
### Security

---
## [1.6.0] - 2025-08-13

### Added
### Changed
- **debt:** Updated to Go `1.25` to leverage the latest language features and performance improvements.

### Deprecated
### Removed
### Fixed
- **defect:** Go `1.25.0`: The `AllocsPerRun` function now panics if parallel tests are running. The result of `AllocsPerRun` is inherently flaky if other tests are running. The new panicking behavior helps catch such bugs.

### Security

---
## [1.5.0] - 2025-07-23

### Added
### Changed
- **risk:** Implemented automatic fork detection and random stream reseeding using process PID tracking. After a process fork, DRBG instances will securely reseed to prevent duplicate random streams in parent and child processes. This eliminates a longstanding CSPRNG risk in forked environments and aligns user-space DRBG safety with that of kernel-backed generators. See `ForkDetectionInterval` in [configuration](../config.go) for tuning and compliance.

### Deprecated
### Removed
### Fixed
### Security

---
## [1.4.0] - 2025-07-20

### Added
### Changed
- **feature:** Prediction resistance mode (`WithPredictionResistance`), enabling NIST SP 800-90A §9.3 compliance by reseeding before every output.
- **feature:** Automatic interval-based reseeding (`WithReseedInterval`) and reseed-on-request-count (`WithReseedRequests`) options.
- **feature:** Added `ReadWithAdditionalInput([]byte)` method for supplying per-call additional input as permitted by the NIST SP 800-90A standard.
- **feature:** Added `Reseed([]byte)` method to manually reseed a DRBG instance with new entropy.

### Deprecated
### Removed
### Fixed
### Security

---
## [1.3.0] - 2025-07-19

### Added
### Changed
- **debt**: Removed unnecessary heap and stack allocations in DRBG fast path
  - Moved all `[16]byte` temporary block buffers to pre-allocated struct fields to eliminate stack allocations and heap escapes.
  - Ensured all block encryption scratch space is reused via DRBG instance buffers rather than stack or ephemeral variables.

### Deprecated
### Removed
### Fixed
### Security

---
## [1.2.0] - 2025-07-19

### Added
### Changed
- **defect:** Move OpenSSF Scorecard badge to Quality section of [README.md](../README.md).
- **risk:** Added fuzz tests.

### Deprecated
### Removed
### Fixed
### Security

---
## [1.1.0] - 2025-07-18

### Added
### Changed
### Deprecated
### Removed
### Fixed
### Security
- **risk:** Added fuzz testing to the `aes-ctr-drbg` module to enhance security and robustness.

---
## [1.0.1] - 2025-07-18

### Added
### Changed
### Deprecated
### Removed
### Fixed
- **defect:** Update `goreleaser` configuration to use the author's username in the changelog.

### Security

---
## [1.0.0] - 2025-07-18

### Added
- **FEATURE:** Initial commit.
### Changed
### Deprecated
### Removed
### Fixed
### Security

[Unreleased]: https://github.com/sixafter/aes-ctr-drbg/compare/v1.21.0...HEAD
[1.21.0]: https://github.com/sixafter/aes-ctr-drbg/compare/v1.20.0...v1.21.0
[1.20.0]: https://github.com/sixafter/aes-ctr-drbg/compare/v1.19.3...v1.20.0
[1.19.3]: https://github.com/sixafter/aes-ctr-drbg/compare/v1.19.2...v1.19.3
[1.19.2]: https://github.com/sixafter/aes-ctr-drbg/compare/v1.19.1...v1.19.2
[1.19.1]: https://github.com/sixafter/aes-ctr-drbg/compare/v1.19.0...v1.19.1
[1.19.0]: https://github.com/sixafter/aes-ctr-drbg/compare/v1.18.0...v1.19.0
[1.18.0]: https://github.com/sixafter/aes-ctr-drbg/compare/v1.17.0...v1.18.0
[1.17.0]: https://github.com/sixafter/aes-ctr-drbg/compare/v1.16.0...v1.17.0
[1.16.0]: https://github.com/sixafter/aes-ctr-drbg/compare/v1.15.0...v1.16.0
[1.15.0]: https://github.com/sixafter/aes-ctr-drbg/compare/v1.14.5...v1.15.0
[1.14.5]: https://github.com/sixafter/aes-ctr-drbg/compare/v1.14.1...v1.14.5
[1.14.1]: https://github.com/sixafter/aes-ctr-drbg/compare/v1.14.0...v1.14.1
[1.14.0]: https://github.com/sixafter/aes-ctr-drbg/compare/v1.13.0...v1.14.0
[1.13.0]: https://github.com/sixafter/aes-ctr-drbg/compare/v1.12.0...v1.13.0
[1.12.0]: https://github.com/sixafter/aes-ctr-drbg/compare/v1.11.0...v1.12.0
[1.11.0]: https://github.com/sixafter/aes-ctr-drbg/compare/v1.10.0...v1.11.0
[1.10.0]: https://github.com/sixafter/aes-ctr-drbg/compare/v1.9.0...v1.10.0
[1.9.0]: https://github.com/sixafter/aes-ctr-drbg/compare/v1.8.0...v1.9.0
[1.8.0]: https://github.com/sixafter/aes-ctr-drbg/compare/v1.7.0...v1.8.0
[1.7.0]: https://github.com/sixafter/aes-ctr-drbg/compare/v1.6.0...v1.7.0
[1.6.0]: https://github.com/sixafter/aes-ctr-drbg/compare/v1.5.0...v1.6.0
[1.5.0]: https://github.com/sixafter/aes-ctr-drbg/compare/v1.4.0...v1.5.0
[1.4.0]: https://github.com/sixafter/aes-ctr-drbg/compare/v1.3.0...v1.4.0
[1.3.0]: https://github.com/sixafter/aes-ctr-drbg/compare/v1.2.0...v1.3.0
[1.2.0]: https://github.com/sixafter/aes-ctr-drbg/compare/v1.1.0...v1.2.0
[1.1.0]: https://github.com/sixafter/aes-ctr-drbg/compare/v1.0.1...v1.1.0
[1.0.1]: https://github.com/sixafter/aes-ctr-drbg/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/sixafter/aes-ctr-drbg/compare/80b4d9e2c5b6a5805bd11741af0eea3d5435889b...v1.0.0

[MUST]: https://datatracker.ietf.org/doc/html/rfc2119
[MUST NOT]: https://datatracker.ietf.org/doc/html/rfc2119
[SHOULD]: https://datatracker.ietf.org/doc/html/rfc2119
[SHOULD NOT]: https://datatracker.ietf.org/doc/html/rfc2119
[MAY]: https://datatracker.ietf.org/doc/html/rfc2119
[SHALL]: https://datatracker.ietf.org/doc/html/rfc2119
[SHALL NOT]: https://datatracker.ietf.org/doc/html/rfc2119
[REQUIRED]: https://datatracker.ietf.org/doc/html/rfc2119
[RECOMMENDED]: https://datatracker.ietf.org/doc/html/rfc2119
[NOT RECOMMENDED]: https://datatracker.ietf.org/doc/html/rfc2119
