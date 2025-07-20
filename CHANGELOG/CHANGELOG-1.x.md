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

[Unreleased]: https://github.com/sixafter/aes-ctr-drbg/compare/v1.3.0...HEAD
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
