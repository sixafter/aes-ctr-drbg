# AES-CTR-DRBG

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

* **Standards-Compliant Implementation**
  Implements NIST SP 800-90A, Revision 1 AES-CTR-DRBG using the Go standard library (`crypto/aes`, `crypto/cipher`). Supports 128-, 192-, and 256-bit keys. State and counter management strictly adhere to the specification.

* **FIPS 140-2 Alignment**
  Designed for use in FIPS 140-2 validated environments and compatible with Go’s FIPS 140 mode (`GODEBUG=fips140=on`). See [FIPS-140.md](FIPS-140.md) for platform guidance.

* **Zero-Allocation Output Path**
  The DRBG is engineered for `0 allocs/op` in its standard `io.Reader` output path, enabling predictable resource usage and high throughput.

* **Asynchronous Key Rotation**
  Supports automatic key rotation after a configurable number of bytes have been generated (`MaxBytesPerKey`). Rekeying occurs asynchronously with exponential backoff and configurable retry limits, reducing long-term key exposure.

* **Prediction Resistance Mode**
  Supports NIST SP 800-90A prediction resistance. When enabled, the DRBG reseeds from system entropy before every output, as required for state compromise resilience.

* **Sharded Pooling for Concurrency**
  Internal state pooling can be sharded across multiple `sync.Pool` instances. The number of shards is configurable, allowing improved performance under concurrent workloads.

* **Extensive Functional Configuration**
  Exposes a comprehensive set of functional options, including:

  * AES key size (128/192/256-bit)
  * Maximum output per key (rekey threshold)
  * Personalization string (domain separation)
  * Shard/pool count
  * Reseed interval and request count
  * Buffer size controls
  * Key rotation and rekey backoff parameters
  * Prediction resistance

* **Thread-Safe and Deterministic**
  All DRBG instances are safe for concurrent use. Output is deterministic for a given seed and personalization.

* **io.Reader Compatibility**
  Implements Go’s `io.Reader` interface for drop-in use as a secure random source.

* **No External Dependencies**
  Depends exclusively on the Go standard library for cryptographic operations.

* **UUID Generation**
  Can be used as a cryptographically secure `io.Reader` with the [`google/uuid`](https://pkg.go.dev/github.com/google/uuid) package and similar libraries.

* **Comprehensive Testing and Fuzzing**
  Includes property-based, fuzz, concurrency, and allocation tests to validate correctness, robustness, and allocation characteristics.

## NIST SP 800-90A Compliance

For a detailed mapping between the implementation and NIST SP 800-90A requirements, see [NIST-SP-800-90A.md](docs/NIST-SP-800-90A.md).

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
BenchmarkDRBG_SyncPool_Baseline_Concurrent/G2-16  	1000000000	         0.6181 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_SyncPool_Baseline_Concurrent/G4-16  	1000000000	         0.5969 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_SyncPool_Baseline_Concurrent/G8-16  	1000000000	         0.5905 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_SyncPool_Baseline_Concurrent/G16-16 	1000000000	         0.5881 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_SyncPool_Baseline_Concurrent/G32-16 	1000000000	         0.5619 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_SyncPool_Baseline_Concurrent/G64-16 	1000000000	         0.5574 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_SyncPool_Baseline_Concurrent/G128-16         	1000000000	         0.5602 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Serial/Serial_Read_16Bytes-16           	44912980	        25.85 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Serial/Serial_Read_32Bytes-16           	39301040	        30.57 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Serial/Serial_Read_64Bytes-16           	29054428	        41.63 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Serial/Serial_Read_256Bytes-16          	10715560	       111.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Serial/Serial_Read_512Bytes-16          	 5944375	       201.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Serial/Serial_Read_4096Bytes-16         	  824330	      1466 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Serial/Serial_Read_16384Bytes-16        	  205989	      5796 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16Bytes_2Goroutines-16         	20332837	        81.01 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16Bytes_4Goroutines-16         	20462421	        83.47 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16Bytes_8Goroutines-16         	20539239	        80.97 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16Bytes_16Goroutines-16        	20763176	        83.57 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16Bytes_32Goroutines-16        	20850964	        82.24 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16Bytes_64Goroutines-16        	20065071	        85.39 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16Bytes_128Goroutines-16       	20253021	        86.33 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_32Bytes_2Goroutines-16         	20839784	        85.07 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_32Bytes_4Goroutines-16         	20984630	        83.24 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_32Bytes_8Goroutines-16         	21078007	        79.86 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_32Bytes_16Goroutines-16        	21062145	        80.23 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_32Bytes_32Goroutines-16        	21413928	        81.68 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_32Bytes_64Goroutines-16        	21627124	        80.56 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_32Bytes_128Goroutines-16       	21518487	        81.18 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_64Bytes_2Goroutines-16         	21210937	        82.78 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_64Bytes_4Goroutines-16         	22693407	        79.86 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_64Bytes_8Goroutines-16         	22735022	        77.76 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_64Bytes_16Goroutines-16        	23076885	        76.72 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_64Bytes_32Goroutines-16        	23874084	        77.75 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_64Bytes_64Goroutines-16        	23542509	        70.23 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_64Bytes_128Goroutines-16       	23997080	        74.13 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_256Bytes_2Goroutines-16        	11627140	       158.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_256Bytes_4Goroutines-16        	11627437	       155.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_256Bytes_8Goroutines-16        	11636860	       149.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_256Bytes_16Goroutines-16       	12079801	       148.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_256Bytes_32Goroutines-16       	12051283	       143.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_256Bytes_64Goroutines-16       	12056252	       135.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_256Bytes_128Goroutines-16      	12118809	       132.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_512Bytes_2Goroutines-16        	11396904	       155.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_512Bytes_4Goroutines-16        	11471233	       155.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_512Bytes_8Goroutines-16        	11518879	       153.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_512Bytes_16Goroutines-16       	11545940	       154.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_512Bytes_32Goroutines-16       	11568248	       152.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_512Bytes_64Goroutines-16       	11636635	       156.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_512Bytes_128Goroutines-16      	11677282	       148.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_4096Bytes_2Goroutines-16       	 7952310	       206.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_4096Bytes_4Goroutines-16       	 7891303	       208.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_4096Bytes_8Goroutines-16       	 7924614	       207.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_4096Bytes_16Goroutines-16      	 7891260	       193.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_4096Bytes_32Goroutines-16      	 7928958	       209.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_4096Bytes_64Goroutines-16      	 7847902	       192.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_4096Bytes_128Goroutines-16     	 7771561	       191.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16384Bytes_2Goroutines-16      	 1729878	       737.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16384Bytes_4Goroutines-16      	 1721493	       754.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16384Bytes_8Goroutines-16      	 1708314	       754.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16384Bytes_16Goroutines-16     	 1674286	       744.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16384Bytes_32Goroutines-16     	 1735957	       756.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16384Bytes_64Goroutines-16     	 1633794	       771.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_Concurrent/Concurrent_Read_16384Bytes_128Goroutines-16    	 1681378	       765.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Sequential/Serial_Read_Large_4096Bytes-16      	  749028	      1490 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Sequential/Serial_Read_Large_16384Bytes-16     	  202723	      5905 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Sequential/Serial_Read_Large_65536Bytes-16     	   51028	     23589 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Sequential/Serial_Read_Large_1048576Bytes-16   	    3130	    380140 ns/op	      15 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_4096Bytes_2Goroutines-16         	 7938176	       181.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_4096Bytes_4Goroutines-16         	 7514082	       193.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_4096Bytes_8Goroutines-16         	 7464860	       188.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_4096Bytes_16Goroutines-16        	 7445499	       189.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_4096Bytes_32Goroutines-16        	 7440232	       199.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_4096Bytes_64Goroutines-16        	 7420635	       189.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_4096Bytes_128Goroutines-16       	 7394299	       194.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_16384Bytes_2Goroutines-16        	 1625565	       795.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_16384Bytes_4Goroutines-16        	 1583449	       754.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_16384Bytes_8Goroutines-16        	 1604343	       753.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_16384Bytes_16Goroutines-16       	 1581939	       779.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_16384Bytes_32Goroutines-16       	 1669182	       778.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_16384Bytes_64Goroutines-16       	 1706570	       744.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_16384Bytes_128Goroutines-16      	 1747254	       773.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_65536Bytes_2Goroutines-16        	  514552	      2652 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_65536Bytes_4Goroutines-16        	  504885	      2699 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_65536Bytes_8Goroutines-16        	  497020	      2635 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_65536Bytes_16Goroutines-16       	  502842	      2673 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_65536Bytes_32Goroutines-16       	  509527	      2684 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_65536Bytes_64Goroutines-16       	  507958	      2632 ns/op	       1 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_65536Bytes_128Goroutines-16      	  509450	      2722 ns/op	       1 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_1048576Bytes_2Goroutines-16      	   30952	     39434 ns/op	       3 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_1048576Bytes_4Goroutines-16      	   31597	     39390 ns/op	       3 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_1048576Bytes_8Goroutines-16      	   31297	     39718 ns/op	       5 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_1048576Bytes_16Goroutines-16     	   30859	     39624 ns/op	       6 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_1048576Bytes_32Goroutines-16     	   31536	     39401 ns/op	       7 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_1048576Bytes_64Goroutines-16     	   31332	     39524 ns/op	      16 B/op	       0 allocs/op
BenchmarkDRBG_Read_LargeSizes_Concurrent/Concurrent_Read_Large_1048576Bytes_128Goroutines-16    	   31790	     39091 ns/op	      20 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes/Serial_Read_Variable_16Bytes-16                                	44671665	        26.25 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes/Serial_Read_Variable_32Bytes-16                                	38319018	        31.44 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes/Serial_Read_Variable_64Bytes-16                                	28390270	        41.87 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes/Serial_Read_Variable_128Bytes-16                               	18427647	        64.69 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes/Serial_Read_Variable_256Bytes-16                               	10597234	       113.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes/Serial_Read_Variable_512Bytes-16                               	 5855329	       229.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes/Serial_Read_Variable_1024Bytes-16                              	 2742954	       387.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes/Serial_Read_Variable_2048Bytes-16                              	 1584697	       742.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes/Serial_Read_Variable_4096Bytes-16                              	  803918	      1489 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_16Bytes_2Goroutines-16     	19970902	        79.48 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_16Bytes_4Goroutines-16     	20621006	        78.10 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_16Bytes_8Goroutines-16     	20840734	        82.56 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_16Bytes_16Goroutines-16    	20802800	        81.93 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_16Bytes_32Goroutines-16    	20894986	        80.47 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_16Bytes_64Goroutines-16    	20891924	        78.08 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_16Bytes_128Goroutines-16   	21197979	        82.07 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_32Bytes_2Goroutines-16     	19203469	        83.56 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_32Bytes_4Goroutines-16     	20180530	        82.40 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_32Bytes_8Goroutines-16     	20800876	        83.30 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_32Bytes_16Goroutines-16    	20901158	        79.16 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_32Bytes_32Goroutines-16    	21472362	        80.99 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_32Bytes_64Goroutines-16    	21482115	        75.88 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_32Bytes_128Goroutines-16   	20990748	        81.06 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_64Bytes_2Goroutines-16     	20392743	        84.37 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_64Bytes_4Goroutines-16     	22488931	        62.82 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_64Bytes_8Goroutines-16     	22626388	        78.22 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_64Bytes_16Goroutines-16    	22863638	        76.92 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_64Bytes_32Goroutines-16    	23986926	        77.12 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_64Bytes_64Goroutines-16    	23878240	        74.41 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_64Bytes_128Goroutines-16   	23789089	        74.53 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_128Bytes_2Goroutines-16    	11371380	       158.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_128Bytes_4Goroutines-16    	11339992	       151.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_128Bytes_8Goroutines-16    	11518612	       154.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_128Bytes_16Goroutines-16   	11495270	       155.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_128Bytes_32Goroutines-16   	11763811	       150.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_128Bytes_64Goroutines-16   	11656874	       107.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_128Bytes_128Goroutines-16  	11795065	       154.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_256Bytes_2Goroutines-16    	11247376	       149.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_256Bytes_4Goroutines-16    	11907961	       154.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_256Bytes_8Goroutines-16    	11763004	       155.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_256Bytes_16Goroutines-16   	11568234	       147.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_256Bytes_32Goroutines-16   	12000434	       148.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_256Bytes_64Goroutines-16   	12056282	       138.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_256Bytes_128Goroutines-16  	12075951	       135.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_512Bytes_2Goroutines-16    	11597779	       142.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_512Bytes_4Goroutines-16    	11727496	       155.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_512Bytes_8Goroutines-16    	11544853	       155.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_512Bytes_16Goroutines-16   	11673974	       149.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_512Bytes_32Goroutines-16   	11749831	       153.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_512Bytes_64Goroutines-16   	11586585	       156.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_512Bytes_128Goroutines-16  	11681427	       152.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_1024Bytes_2Goroutines-16   	12218305	       148.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_1024Bytes_4Goroutines-16   	12357247	       148.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_1024Bytes_8Goroutines-16   	12306213	       150.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_1024Bytes_16Goroutines-16  	12367579	       149.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_1024Bytes_32Goroutines-16  	12521710	       146.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_1024Bytes_64Goroutines-16  	12427521	       143.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_1024Bytes_128Goroutines-16 	12349743	       146.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_2048Bytes_2Goroutines-16   	10496486	       117.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_2048Bytes_4Goroutines-16   	12406910	       120.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_2048Bytes_8Goroutines-16   	12478713	       126.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_2048Bytes_16Goroutines-16  	12387714	       113.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_2048Bytes_32Goroutines-16  	12250074	       110.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_2048Bytes_64Goroutines-16  	12214362	       118.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_2048Bytes_128Goroutines-16 	12112784	       113.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_4096Bytes_2Goroutines-16   	 7166356	       193.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_4096Bytes_4Goroutines-16   	 7244019	       195.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_4096Bytes_8Goroutines-16   	 7225875	       211.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_4096Bytes_16Goroutines-16  	 7206831	       206.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_4096Bytes_32Goroutines-16  	 7185728	       202.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_4096Bytes_64Goroutines-16  	 7153624	       204.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_VariableSizes_Concurrent/Concurrent_Read_Variable_4096Bytes_128Goroutines-16 	 7159110	       198.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Serial_Read_Extreme_10485760Bytes-16                            	     302	   4373806 ns/op	     163 B/op	       0 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_2Goroutines-16            	    3698	    360310 ns/op	      78 B/op	       0 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_4Goroutines-16            	    3171	    365370 ns/op	     102 B/op	       0 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_8Goroutines-16            	    3140	    366112 ns/op	     114 B/op	       0 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_16Goroutines-16           	    3188	    370715 ns/op	     120 B/op	       0 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_32Goroutines-16           	    2833	    371710 ns/op	     188 B/op	       1 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_64Goroutines-16           	    3078	    373323 ns/op	     198 B/op	       1 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_10485760Bytes_128Goroutines-16          	    3031	    373651 ns/op	     292 B/op	       2 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Serial_Read_Extreme_52428800Bytes-16                            	      56	  19282917 ns/op	     882 B/op	       1 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_2Goroutines-16            	     669	   1805320 ns/op	     320 B/op	       1 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_4Goroutines-16            	     612	   1829114 ns/op	     398 B/op	       2 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_8Goroutines-16            	     615	   1848068 ns/op	     453 B/op	       2 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_16Goroutines-16           	     622	   1848697 ns/op	     454 B/op	       3 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_32Goroutines-16           	     625	   1846224 ns/op	     572 B/op	       4 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_64Goroutines-16           	     592	   1872232 ns/op	     732 B/op	       6 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_52428800Bytes_128Goroutines-16          	     607	   1875778 ns/op	    1031 B/op	      10 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Serial_Read_Extreme_104857600Bytes-16                           	      30	  38446665 ns/op	    1237 B/op	       2 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_2Goroutines-16           	     319	   3647249 ns/op	     549 B/op	       2 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_4Goroutines-16           	     301	   3692198 ns/op	     617 B/op	       3 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_8Goroutines-16           	     310	   3658588 ns/op	     746 B/op	       4 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_16Goroutines-16          	     316	   3711240 ns/op	     837 B/op	       5 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_32Goroutines-16          	     319	   3721948 ns/op	    1058 B/op	       7 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_64Goroutines-16          	     308	   3730344 ns/op	    1203 B/op	      11 allocs/op
BenchmarkDRBG_Read_ExtremeSizes/Concurrent_Read_Extreme_104857600Bytes_128Goroutines-16         	     315	   3699707 ns/op	    1457 B/op	      17 allocs/op
BenchmarkDRBG_Read_WithKeyRotation-16                                                           	 5226307	       228.5 ns/op	     192 B/op	       1 allocs/op
BenchmarkDRBG_Read_PredictionResistance-16                                                      	 2630257	       455.1 ns/op	     634 B/op	       3 allocs/op
PASS
ok  	github.com/sixafter/aes-ctr-drbg	313.461s
```

</details>

### UUID Generation with Google UUID and ctrdrbg

Here's a summary of the benchmark results comparing the default random reader for Google's [UUID](https://pkg.go.dev/github.com/google/uuid) package and the ctrdrbg-based UUID generation:

| Benchmark Scenario                  | Default ns/op | CTRDRBG ns/op | % Faster (ns/op) | Default B/op | CTRDRBG B/op | Default allocs/op | CTRDRBG allocs/op |
|-------------------------------------|--------------:|--------------:|-----------------:|-------------:|-------------:|------------------:|------------------:|
| v4 Serial                           |        180.6  |        40.78  |         77.4%    |         16   |         16   |                1  |                1  |
| v4 Parallel                         |        445.4  |        10.56  |         97.6%    |         16   |         16   |                1  |                1  |
| v4 Concurrent (2 goroutines)        |        413.2  |        21.91  |         94.7%    |         16   |         16   |                1  |                1  |
| v4 Concurrent (4 goroutines)        |        428.5  |        12.77  |         97.0%    |         16   |         16   |                1  |                1  |
| v4 Concurrent (8 goroutines)        |        484.6  |         9.74  |         98.0%    |         16   |         16   |                1  |                1  |
| v4 Concurrent (16 goroutines)       |        458.2  |         7.67  |         98.3%    |         16   |         16   |                1  |                1  |
| v4 Concurrent (32 goroutines)       |        506.3  |         7.69  |         98.5%    |         16   |         16   |                1  |                1  |
| v4 Concurrent (64 goroutines)       |        506.9  |         7.64  |         98.5%    |         16   |         16   |                1  |                1  |
| v4 Concurrent (128 goroutines)      |        508.2  |         7.63  |         98.5%    |         16   |         16   |                1  |                1  |
| v4 Concurrent (256 goroutines)      |        511.8  |         7.79  |         98.5%    |         16   |         16   |                1  |                1  |

Notes:
- “Default” refers to the baseline Go `crypto/rand` source.
- “CTRDRBG” refers to your AES-CTR-DRBG implementation.
- “% Faster (ns/op)” is computed as `100 * (Default - CTRDRBG) / Default`, rounded.

<details>
  <summary>Expand to see results</summary>

  ```shell
make bench-uuid
go test -bench='^BenchmarkUUID_' -run=^$ -benchmem -memprofile=mem.out -cpuprofile=cpu.out .
goos: darwin
goarch: arm64
pkg: github.com/sixafter/aes-ctr-drbg
cpu: Apple M4 Max
BenchmarkUUID_v4_Default_Serial-16        	 6473760	       180.6 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Parallel-16      	 2705866	       445.4 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_2-16         	 2883284	       413.2 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_4-16         	 2806682	       428.5 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_8-16         	 2462146	       484.6 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_16-16        	 2685201	       458.2 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_32-16        	 2366074	       506.3 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_64-16        	 2358429	       506.9 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_128-16       	 2388648	       508.2 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_Default_Concurrent/Goroutines_256-16       	 2364384	       511.8 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CTRDRBG_Serial-16                          	29120706	        40.78 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CTRDRBG_Parallel-16                        	100000000	        10.56 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CTRDRBG_Concurrent/Goroutines_2-16         	52686843	        21.91 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CTRDRBG_Concurrent/Goroutines_4-16         	92968908	        12.77 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CTRDRBG_Concurrent/Goroutines_8-16         	121979662	         9.741 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CTRDRBG_Concurrent/Goroutines_16-16        	153623710	         7.668 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CTRDRBG_Concurrent/Goroutines_32-16        	154797238	         7.688 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CTRDRBG_Concurrent/Goroutines_64-16        	156757164	         7.641 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CTRDRBG_Concurrent/Goroutines_128-16       	156462766	         7.632 ns/op	      16 B/op	       1 allocs/op
BenchmarkUUID_v4_CTRDRBG_Concurrent/Goroutines_256-16       	154197008	         7.795 ns/op	      16 B/op	       1 allocs/op
PASS
ok  	github.com/sixafter/aes-ctr-drbg	33.515s
  ```
</details>

---

## Contributing

Contributions are welcome. See [CONTRIBUTING](CONTRIBUTING.md)

---

## License

This project is licensed under the [Apache 2.0 License](https://choosealicense.com/licenses/apache-2.0/). See [LICENSE](LICENSE) file.
