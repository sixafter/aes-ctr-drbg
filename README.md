# A Deterministic Random Bit Generator based on AES in Counter mode (AES-CTR-DRBG) as specified in NIST SP 800-90A

[![Go Report Card](https://goreportcard.com/badge/github.com/sixafter/aes-ctr-drbg)](https://goreportcard.com/report/github.com/sixafter/aes-ctr-drbg)
[![License: Apache 2.0](https://img.shields.io/badge/license-Apache%202.0-blue?style=flat-square)](LICENSE)
[![Go](https://img.shields.io/github/go-mod/go-version/sixafter/aes-ctr-drbg)](https://img.shields.io/github/go-mod/go-version/sixafter/aes-ctr-drbg)
[![Go Reference](https://pkg.go.dev/badge/github.com/sixafter/aes-ctr-drbg.svg)](https://pkg.go.dev/github.com/sixafter/aes-ctr-drbg)
[![FIPS‑140 Mode Compatible](https://img.shields.io/badge/FIPS‑140--Mode-Compatible-brightgreen)](FIPS‑140.md)

---

## Status

### Build & Test

[![CI](https://github.com/sixafter/aes-ctr-drbg/workflows/ci/badge.svg)](https://github.com/sixafter/aes-ctr-drbg/actions)
[![GitHub issues](https://img.shields.io/github/issues/sixafter/aes-ctr-drbg)](https://github.com/sixafter/aes-ctr-drbg/issues)

### Quality

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=six-after_aes-ctr-drbg&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=six-after_aes-ctr-drbg)
![CodeQL](https://github.com/sixafter/aes-ctr-drbg/actions/workflows/codeql-analysis.yaml/badge.svg)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=six-after_aes-ctr-drbg&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=six-after_aes-ctr-drbg)
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/sixafter/aes-ctr-drbg/badge)](https://scorecard.dev/viewer/?uri=github.com/sixafter/aes-ctr-drbg)

### Package and Deploy

[![Release](https://github.com/sixafter/aes-ctr-drbg/workflows/release/badge.svg)](https://github.com/sixafter/aes-ctr-drbg/actions)

---
## Overview 

AES-CTR-DRBG (Deterministic Random Bit Generator based on AES in Counter mode) is a cryptographically secure pseudo-random number generator (CSPRNG) defined by [NIST SP 800-90A Rev. 1](https://csrc.nist.gov/pubs/sp/800/90/a/r1/final). It's widely used in high-assurance systems, including those requiring [FIPS 140-2](https://csrc.nist.gov/pubs/fips/140-2/upd2/final) compliance.

AES-CTR-DRBG is designed for environments requiring deterministic, reproducible, and FIPS‑140-compatible random bit generation. This module is suitable for any application that needs strong cryptographic assurance or must comply with regulated environments (e.g., FedRAMP, FIPS, PCI, HIPAA).

The module uses only Go standard library crypto primitives (`crypto/aes` and `crypto/cipher`), making it safe for use in FIPS 140-validated Go runtimes. No third-party, homegrown, or experimental ciphers are used.

Please see the [godoc](https://pkg.go.dev/github.com/sixafter/aes-ctr-drbg) for detailed documentation.

---

## FIPS‑140 Mode

See [FIPS‑140.md](FIPS-140.md) for compliance, deployment, and configuration guidance.

---

## Features

- **NIST SP 800-90A AES-CTR-DRBG Implementation:** Implements the Deterministic Random Bit Generator (DRBG) construction defined in [NIST SP 800-90A, Revision 1](https://csrc.nist.gov/pubs/sp/800/90/a/r1/final), using the AES block cipher in counter (CTR) mode.
  - Supports 128-, 192-, and 256-bit AES keys, with correct counter and state management as specified by the standard.
- **FIPS 140-2 Alignment:** Designed to operate in FIPS 140-2 validated environments and compatible with Go’s FIPS 140 mode (`GODEBUG=fips140=on`).
  - Uses only cryptographic primitives from the Go standard library.
  - For platform-specific guidance and deployment instructions, see [FIPS‑140.md](FIPS-140.md).
- **Optimized for Low Allocations:** Carefully structured to minimize heap allocations, reducing memory overhead and improving cache locality. This optimization is crucial for applications where performance and resource usage are critical.
  - `0 allocs/op` for `io.Reader` interface
- **Stateless and Concurrent Operation:** Each DRBG instance is safe for concurrent use and fully encapsulates its cryptographic state. 
  - The design supports independent operation across multiple instances, enabling scalable use in high-concurrency environments.
- **Configurable Entropy and Personalization:** Accepts externally supplied entropy sources and personalization strings, enabling domain separation and deterministic output for compliance with best practices and advanced use cases.
- **io.Reader Compatibility:** Fully satisfies Go’s `io.Reader` interface, allowing seamless integration with packages and APIs expecting a secure random source.
- **No External Dependencies:** Depends solely on the Go standard library, ensuring a lightweight and portable implementation.
- **UUID Generation Source:** Can be used as the `io.Reader` source for UUID generation with the [`google/uuid`](https://pkg.go.dev/github.com/google/uuid) package and similar libraries, providing cryptographically secure, deterministic UUIDs using AES-CTR-DRBG.

---

## Verify with Cosign

[Cosign](https://github.com/sigstore/cosign) is used to sign releases for integrity verification.

To verify the integrity of the release tarball, you can use Cosign to check the signature and checksums. Follow these steps:

```sh
# Fetch the latest release tag from GitHub API (e.g., "v1.3.0")
TAG=$(curl -s https://api.github.com/repos/sixafter/aes-ctr-drbg/releases/latest | jq -r .tag_name)

# Remove leading "v" for filenames (e.g., "v1.3.0" -> "1.3.0")
VERSION=${TAG#v}

# Verify the release tarball
cosign verify-blob \
  --key https://raw.githubusercontent.com/sixafter/aes-ctr-drbg/main/cosign.pub \
  --signature aes-ctr-drbg-${VERSION}.tar.gz.sig \
  aes-ctr-drbg-${VERSION}.tar.gz

# Download checksums.txt and its signature from the latest release assets
curl -LO https://github.com/sixafter/aes-ctr-drbg/releases/download/${TAG}/checksums.txt
curl -LO https://github.com/sixafter/aes-ctr-drbg/releases/download/${TAG}/checksums.txt.sig

# Verify checksums.txt with cosign
cosign verify-blob \
  --key https://raw.githubusercontent.com/sixafter/aes-ctr-drbg/main/cosign.pub \
  --signature checksums.txt.sig \
  checksums.txt
```

If valid, Cosign will output:

```shell
Verified OK
```

---

## Installation

```bash
go get -u github.com/sixafter/aes-ctr-drbg
```

---

## Usage

### Basic Usage: Generate Secure Random Bytes With Reader

```go
package main

import (
	"fmt"
	"log"

	"github.com/sixafter/aes-ctr-drbg"
)

func main() {
	buf := make([]byte, 64)
	n, err := ctrdrbg.Reader.Read(buf)
	if err != nil {
		log.Fatalf("failed to read random bytes: %v", err)
	}
	fmt.Printf("Read %d random bytes: %x\n", n, buf)
}
```

### Basic Usage: Generate Secure Random Bytes with NewReader

```go
package main

import (
	"fmt"
	"log"

	"github.com/sixafter/aes-ctr-drbg"
)

func main() {
	// Example: AES-256 (32 bytes) key
	r, err := ctrdrbg.NewReader(ctrdrbg.WithKeySize(ctrdrbg.KeySize256))
	if err != nil {
		log.Fatalf("failed to create ctrdrbg.Reader: %v", err)
	}

	buf := make([]byte, 64)
	n, err := r.Read(buf)
	if err != nil {
		log.Fatalf("failed to read random bytes: %v", err)
	}
	fmt.Printf("Read %d random bytes: %x\n", n, buf)
}
```

### Using Personalization and Additional Input

```go
package main

import (
	"fmt"
	"log"

	"github.com/sixafter/aes-ctr-drbg"
)

func main() {
	r, err := ctrdrbg.NewReader(
		ctrdrbg.WithPersonalization([]byte("service-id-1")),
		ctrdrbg.WithKeySize(ctrdrbg.KeySize256), // AES-256
	)
	if err != nil {
		log.Fatalf("failed to create ctrdrbg.Reader: %v", err)
	}

	buf := make([]byte, 64)
	n, err := r.Read(buf)
	if err != nil {
		log.Fatalf("failed to read random bytes: %v", err)
	}
	fmt.Printf("Read %d random bytes: %x\n", n, buf)
}
```

---

## Performance Benchmarks

### Raw Random Byte Generation

These `ctrdrbg.Reader` benchmarks demonstrate the package's focus on minimizing latency, memory usage, and allocation overhead, making it suitable for high-performance applications.

<details>
  <summary>Expand to see results</summary>

```shell
make bench
go test -bench='^BenchmarkDRBG_' -run=^$ -benchmem -memprofile=mem.out -cpuprofile=cpu.out .
goos: darwin
goarch: arm64
pkg: github.com/sixafter/aes-ctr-drbg
cpu: Apple M4 Max
BenchmarkDRBG_SyncPool_Baseline_Concurrent/G2-16  	1000000000	         0.6077 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_SyncPool_Baseline_Concurrent/G4-16  	1000000000	         0.6087 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_SyncPool_Baseline_Concurrent/G8-16  	1000000000	         0.5845 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_SyncPool_Baseline_Concurrent/G16-16 	1000000000	         0.6001 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_SyncPool_Baseline_Concurrent/G32-16 	1000000000	         0.5659 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_SyncPool_Baseline_Concurrent/G64-16 	1000000000	         0.5601 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_SyncPool_Baseline_Concurrent/G128-16         	1000000000	         0.5605 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Serial/Serial_Read_16Bytes-16           	46836044	        25.03 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Serial/Serial_Read_32Bytes-16           	39548868	        29.92 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Serial/Serial_Read_64Bytes-16           	28945945	        41.07 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Serial/Serial_Read_256Bytes-16          	11048659	       110.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Serial/Serial_Read_512Bytes-16          	 6105094	       196.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Serial/Serial_Read_4096Bytes-16         	  842362	      1453 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Serial/Serial_Read_16384Bytes-16        	  211924	      5640 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16Bytes_2Goroutines-16         	19675221	        84.26 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16Bytes_4Goroutines-16         	20334344	        84.37 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16Bytes_8Goroutines-16         	20486322	        83.76 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16Bytes_16Goroutines-16        	20606031	        84.26 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16Bytes_32Goroutines-16        	20096574	        85.28 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16Bytes_64Goroutines-16        	20174424	        85.55 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16Bytes_128Goroutines-16       	20311455	        83.34 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_32Bytes_2Goroutines-16         	20105216	        84.45 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_32Bytes_4Goroutines-16         	21605125	        84.04 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_32Bytes_8Goroutines-16         	22053130	        83.30 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_32Bytes_16Goroutines-16        	21417176	        81.00 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_32Bytes_32Goroutines-16        	22384962	        79.62 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_32Bytes_64Goroutines-16        	22610188	        77.46 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_32Bytes_128Goroutines-16       	22689813	        77.92 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_64Bytes_2Goroutines-16         	22438783	        83.34 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_64Bytes_4Goroutines-16         	22370060	        80.27 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_64Bytes_8Goroutines-16         	23965848	        79.43 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_64Bytes_16Goroutines-16        	24108346	        78.08 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_64Bytes_32Goroutines-16        	24453572	        73.57 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_64Bytes_64Goroutines-16        	24086893	        70.55 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_64Bytes_128Goroutines-16       	25006185	        51.41 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_256Bytes_2Goroutines-16        	11825310	       123.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_256Bytes_4Goroutines-16        	11816847	       149.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_256Bytes_8Goroutines-16        	11843376	       150.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_256Bytes_16Goroutines-16       	11676459	       142.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_256Bytes_32Goroutines-16       	12061593	       147.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_256Bytes_64Goroutines-16       	11899906	       139.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_256Bytes_128Goroutines-16      	12015664	       136.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_512Bytes_2Goroutines-16        	11359638	       157.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_512Bytes_4Goroutines-16        	11516784	       156.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_512Bytes_8Goroutines-16        	11563575	       155.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_512Bytes_16Goroutines-16       	11656911	       154.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_512Bytes_32Goroutines-16       	11668029	       155.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_512Bytes_64Goroutines-16       	11694646	       152.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_512Bytes_128Goroutines-16      	11666110	       154.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_4096Bytes_2Goroutines-16       	 7994008	       219.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_4096Bytes_4Goroutines-16       	 7988438	       211.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_4096Bytes_8Goroutines-16       	 7964991	       212.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_4096Bytes_16Goroutines-16      	 7972845	       215.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_4096Bytes_32Goroutines-16      	 7893076	       211.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_4096Bytes_64Goroutines-16      	 7924053	       216.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_4096Bytes_128Goroutines-16     	 7816192	       213.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16384Bytes_2Goroutines-16      	 1693452	       822.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16384Bytes_4Goroutines-16      	 1630606	       711.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16384Bytes_8Goroutines-16      	 1564152	       827.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16384Bytes_16Goroutines-16     	 1606071	       833.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16384Bytes_32Goroutines-16     	 1553371	       834.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16384Bytes_64Goroutines-16     	 1542376	       820.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16384Bytes_128Goroutines-16    	 1571803	       835.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Sequential/Serial_Read_Large_4096Bytes-16      	  780607	      1445 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Sequential/Serial_Read_Large_16384Bytes-16     	  209959	      5699 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Sequential/Serial_Read_Large_65536Bytes-16     	   51926	     23134 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Sequential/Serial_Read_Large_1048576Bytes-16   	    3190	    374031 ns/op	      15 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_4096Bytes_2Goroutines-16         	 7994466	       219.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_4096Bytes_4Goroutines-16         	 7573742	       219.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_4096Bytes_8Goroutines-16         	 7451446	       221.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_4096Bytes_16Goroutines-16        	 6927225	       210.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_4096Bytes_32Goroutines-16        	 7471450	       218.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_4096Bytes_64Goroutines-16        	 7497517	       219.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_4096Bytes_128Goroutines-16       	 7443847	       219.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_16384Bytes_2Goroutines-16        	 1571527	       840.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_16384Bytes_4Goroutines-16        	 1590352	       840.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_16384Bytes_8Goroutines-16        	 1561832	       840.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_16384Bytes_16Goroutines-16       	 1593702	       840.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_16384Bytes_32Goroutines-16       	 1545817	       839.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_16384Bytes_64Goroutines-16       	 1569508	       822.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_16384Bytes_128Goroutines-16      	 1584163	       842.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_65536Bytes_2Goroutines-16        	  500068	      2769 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_65536Bytes_4Goroutines-16        	  510716	      2768 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_65536Bytes_8Goroutines-16        	  520544	      2748 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_65536Bytes_16Goroutines-16       	  511690	      2735 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_65536Bytes_32Goroutines-16       	  527379	      2758 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_65536Bytes_64Goroutines-16       	  521968	      2764 ns/op	       1 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_65536Bytes_128Goroutines-16      	  527312	      2742 ns/op	       1 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_1048576Bytes_2Goroutines-16      	   31135	     39914 ns/op	       3 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_1048576Bytes_4Goroutines-16      	   31465	     39595 ns/op	       3 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_1048576Bytes_8Goroutines-16      	   31364	     39165 ns/op	       4 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_1048576Bytes_16Goroutines-16     	   31669	     39074 ns/op	       4 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_1048576Bytes_32Goroutines-16     	   31144	     38845 ns/op	       6 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_1048576Bytes_64Goroutines-16     	   31347	     38644 ns/op	      16 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_1048576Bytes_128Goroutines-16    	   31270	     38876 ns/op	      18 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes/Serial_Read_Variable_16Bytes-16                                	46493029	        25.32 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes/Serial_Read_Variable_32Bytes-16                                	39151179	        30.50 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes/Serial_Read_Variable_64Bytes-16                                	28405560	        41.98 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes/Serial_Read_Variable_128Bytes-16                               	18663969	        64.03 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes/Serial_Read_Variable_256Bytes-16                               	10717100	       110.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes/Serial_Read_Variable_512Bytes-16                               	 6030423	       198.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes/Serial_Read_Variable_1024Bytes-16                              	 3179872	       375.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes/Serial_Read_Variable_2048Bytes-16                              	 1639771	       732.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes/Serial_Read_Variable_4096Bytes-16                              	  812773	      1442 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_16Bytes_2Goroutines-16     	20376266	        84.23 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_16Bytes_4Goroutines-16     	20499051	        81.80 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_16Bytes_8Goroutines-16     	20043440	        80.11 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_16Bytes_16Goroutines-16    	20619825	        83.98 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_16Bytes_32Goroutines-16    	20788114	        84.20 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_16Bytes_64Goroutines-16    	20955132	        81.88 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_16Bytes_128Goroutines-16   	20886955	        83.14 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_32Bytes_2Goroutines-16     	20132496	        86.21 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_32Bytes_4Goroutines-16     	21361308	        83.56 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_32Bytes_8Goroutines-16     	21408373	        81.61 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_32Bytes_16Goroutines-16    	22204917	        83.02 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_32Bytes_32Goroutines-16    	22034689	        79.61 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_32Bytes_64Goroutines-16    	22582956	        79.14 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_32Bytes_128Goroutines-16   	22777114	        77.90 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_64Bytes_2Goroutines-16     	22852154	        81.61 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_64Bytes_4Goroutines-16     	24239893	        81.13 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_64Bytes_8Goroutines-16     	24178334	        78.63 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_64Bytes_16Goroutines-16    	23951278	        79.39 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_64Bytes_32Goroutines-16    	24378289	        76.62 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_64Bytes_64Goroutines-16    	24630140	        71.35 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_64Bytes_128Goroutines-16   	24671052	        71.60 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_128Bytes_2Goroutines-16    	11239616	       155.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_128Bytes_4Goroutines-16    	11660102	       156.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_128Bytes_8Goroutines-16    	11552350	       156.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_128Bytes_16Goroutines-16   	11690331	       148.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_128Bytes_32Goroutines-16   	11779237	       152.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_128Bytes_64Goroutines-16   	11698827	       141.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_128Bytes_128Goroutines-16  	11797335	       151.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_256Bytes_2Goroutines-16    	11392652	       157.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_256Bytes_4Goroutines-16    	11828292	       154.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_256Bytes_8Goroutines-16    	11691841	       154.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_256Bytes_16Goroutines-16   	11759618	       154.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_256Bytes_32Goroutines-16   	12011761	       146.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_256Bytes_64Goroutines-16   	12179396	       140.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_256Bytes_128Goroutines-16  	12070298	       137.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_512Bytes_2Goroutines-16    	11552596	       144.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_512Bytes_4Goroutines-16    	11572473	       155.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_512Bytes_8Goroutines-16    	11563664	       155.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_512Bytes_16Goroutines-16   	11580985	       155.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_512Bytes_32Goroutines-16   	11557320	       155.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_512Bytes_64Goroutines-16   	11634788	       154.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_512Bytes_128Goroutines-16  	11663238	       155.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_1024Bytes_2Goroutines-16   	12286054	       151.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_1024Bytes_4Goroutines-16   	12339008	       149.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_1024Bytes_8Goroutines-16   	12326107	       152.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_1024Bytes_16Goroutines-16  	12362986	       143.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_1024Bytes_32Goroutines-16  	12432188	       147.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_1024Bytes_64Goroutines-16  	12478946	       151.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_1024Bytes_128Goroutines-16 	12471019	       148.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_2048Bytes_2Goroutines-16   	12544668	       129.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_2048Bytes_4Goroutines-16   	10214352	       126.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_2048Bytes_8Goroutines-16   	12378098	       128.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_2048Bytes_16Goroutines-16  	12476562	       128.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_2048Bytes_32Goroutines-16  	12464899	       125.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_2048Bytes_64Goroutines-16  	12445813	       122.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_2048Bytes_128Goroutines-16 	12237784	       123.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_4096Bytes_2Goroutines-16   	 7323421	       220.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_4096Bytes_4Goroutines-16   	 7195437	       223.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_4096Bytes_8Goroutines-16   	 7296626	       257.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_4096Bytes_16Goroutines-16  	 7177464	       220.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_4096Bytes_32Goroutines-16  	 7279090	       191.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_4096Bytes_64Goroutines-16  	 7355392	       207.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_4096Bytes_128Goroutines-16 	 7223954	       223.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Serial_Read_Extreme_10485760Bytes-16                            	     309	   3758409 ns/op	     158 B/op	       0 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_2Goroutines-16            	    3723	    359122 ns/op	      80 B/op	       0 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_4Goroutines-16            	    3214	    364965 ns/op	      98 B/op	       0 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_8Goroutines-16            	    3111	    365863 ns/op	     108 B/op	       0 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_16Goroutines-16           	    3224	    365678 ns/op	     119 B/op	       0 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_32Goroutines-16           	    3201	    367356 ns/op	     123 B/op	       0 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_64Goroutines-16           	    3127	    369390 ns/op	     151 B/op	       1 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_128Goroutines-16          	    3172	    372482 ns/op	     271 B/op	       2 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Serial_Read_Extreme_52428800Bytes-16                            	      56	  19236001 ns/op	     873 B/op	       1 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_2Goroutines-16            	     686	   1824134 ns/op	     305 B/op	       1 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_4Goroutines-16            	     604	   1848823 ns/op	     384 B/op	       2 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_8Goroutines-16            	     597	   1815822 ns/op	     424 B/op	       2 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_16Goroutines-16           	     632	   1831721 ns/op	     408 B/op	       2 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_32Goroutines-16           	     633	   1834422 ns/op	     730 B/op	       5 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_64Goroutines-16           	     594	   1834792 ns/op	     825 B/op	       7 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_128Goroutines-16          	     640	   1868634 ns/op	     977 B/op	      11 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Serial_Read_Extreme_104857600Bytes-16                           	      30	  37860900 ns/op	    1426 B/op	       2 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_2Goroutines-16           	     343	   3610056 ns/op	     512 B/op	       2 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_4Goroutines-16           	     308	   3619210 ns/op	     628 B/op	       3 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_8Goroutines-16           	     282	   3716647 ns/op	     784 B/op	       4 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_16Goroutines-16          	     328	   3694658 ns/op	     782 B/op	       5 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_32Goroutines-16          	     314	   3649463 ns/op	    1059 B/op	       7 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_64Goroutines-16          	     332	   3691214 ns/op	    1146 B/op	      10 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_128Goroutines-16         	     325	   3657829 ns/op	    1406 B/op	      17 allocs/op
PASS
ok  	github.com/sixafter/aes-ctr-drbg	319.559s
```

</details>

### UUID Generation with Google UUID and ctrdrbg

Here's a summary of the benchmark results comparing the default random reader for Google's [UUID](https://pkg.go.dev/github.com/google/uuid) package and the ctrdrbg-based UUID generation:

| Benchmark Scenario                         | Default ns/op | CTRDRBG ns/op | % Faster (ns/op) | Default B/op | CTRDRBG B/op | Default allocs/op | CTRDRBG allocs/op |
|--------------------------------------------|---------------:|---------------:|------------------:|--------------:|--------------:|-------------------:|-------------------:|
| v4 Serial                                   |         184.1  |         41.56 |           77.4%   |          16   |          16   |                 1  |                 1  |
| v4 Parallel                                 |         457.3  |         11.25 |           97.5%   |          16   |          16   |                 1  |                 1  |
| v4 Concurrent (2 goroutines)                |         375.5  |         22.41 |           94.0%   |          16   |          16   |                 1  |                 1  |
| v4 Concurrent (4 goroutines)                |         462.7  |         12.97 |           97.2%   |          16   |          16   |                 1  |                 1  |
| v4 Concurrent (8 goroutines)                |         485.7  |          9.92 |           98.0%   |          16   |          16   |                 1  |                 1  |
| v4 Concurrent (16 goroutines)               |         446.7  |          7.95 |           98.2%   |          16   |          16   |                 1  |                 1  |
| v4 Concurrent (32 goroutines)               |         509.6  |          8.04 |           98.4%   |          16   |          16   |                 1  |                 1  |
| v4 Concurrent (64 goroutines)               |         524.6  |          7.97 |           98.5%   |          16   |          16   |                 1  |                 1  |
| v4 Concurrent (128 goroutines)              |         515.4  |          7.96 |           98.5%   |          16   |          16   |                 1  |                 1  |
| v4 Concurrent (256 goroutines)              |         508.7  |          8.15 |           98.4%   |          16   |          16   |                 1  |                 1  |

<details>
  <summary>Expand to see results</summary>

  ```shell
make bench-uuid
go test -bench='^BenchmarkUUID_' -run=^$ -benchmem -memprofile=mem.out -cpuprofile=cpu.out .
goos: darwin
goarch: arm64
pkg: github.com/sixafter/aes-ctr-drbg
cpu: Apple M4 Max
BenchmarkUUID_v4_Default_Serial-16        	 6262804	       184.1 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Parallel-16      	 2643778	       457.3 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_2-16         	 3218464	       375.5 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_4-16         	 2626228	       462.7 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_8-16         	 2488254	       485.7 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_16-16        	 2683616	       446.7 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_32-16        	 2334570	       509.6 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_64-16        	 2334774	       524.6 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_128-16       	 2364202	       515.4 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_256-16       	 2315127	       508.7 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CTRDRBG_Serial-16                          	27666423	        41.56 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CTRDRBG_Parallel-16                        	100000000	        11.25 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CTRDRBG_Concurrent/Goroutines_2-16         	51704079	        22.41 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CTRDRBG_Concurrent/Goroutines_4-16         	92681390	        12.97 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CTRDRBG_Concurrent/Goroutines_8-16         	120022134	         9.919 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CTRDRBG_Concurrent/Goroutines_16-16        	148474417	         7.946 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CTRDRBG_Concurrent/Goroutines_32-16        	150344217	         8.039 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CTRDRBG_Concurrent/Goroutines_64-16        	149777127	         7.969 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CTRDRBG_Concurrent/Goroutines_128-16       	150662916	         7.964 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CTRDRBG_Concurrent/Goroutines_256-16       	148184919	         8.146 ns/op	      16 B/op	       1 allocs/op
PASS
ok  	github.com/sixafter/aes-ctr-drbg	33.722s
  ```
</details>

---

## Contributing

Contributions are welcome. See [CONTRIBUTING](CONTRIBUTING.md)

---

## License

This project is licensed under the [Apache 2.0 License](https://choosealicense.com/licenses/apache-2.0/). See [LICENSE](LICENSE) file.
