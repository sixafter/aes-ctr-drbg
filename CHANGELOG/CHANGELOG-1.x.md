# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Date format: `YYYY-MM-DD`

---
## [Unreleased]

### Added
### Changed
- **debt:** Modified README to include examples of how AES-CTR-DRBG is used by [NanoID](https://github.com/sixafter/nanoid) when FIPS mode is enabled.

### Deprecated
### Removed
### Fixed
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
- **feature:** Prediction resistance mode (`WithPredictionResistance`), enabling NIST SP 800-90A Section 9.3 compliance by reseeding before every output.
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

[Unreleased]: https://github.com/sixafter/aes-ctr-drbg/compare/v1.14.1...HEAD
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
