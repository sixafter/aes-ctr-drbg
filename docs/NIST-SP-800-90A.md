# NIST SP 800-90A CTR_DRBG: Specification Mapping

This package implements the CTR_DRBG mechanism of [NIST SP 800-90A Rev. 1](https://csrc.nist.gov/pubs/sp/800/90/a/r1/final)
(§10.2.1) with AES-128, AES-192, or AES-256 and no derivation function.

Conformance is opt-in: enable it with `WithBacktrackingResistance`. The default mode generates output the same way,
but derives its initial state directly from entropy, replaces its state on reseed, and does not update the key and
counter after each request; it performs 0 allocations per `Read` call. With `WithBacktrackingResistance`, a `Read`
call performs 2 allocations.

The conforming mode is tested against the NIST CAVP CTR_DRBG test vectors (`drbgtestvectors.zip`: `no_reseed`,
`pr_false`, and `pr_true`, "no df" sections for each key size). The tests named below embed a subset of those vectors.

| Requirement | Default mode | With `WithBacktrackingResistance` | Tests |
|---|---|---|---|
| **Instantiate** (§10.2.1.3.1) | Key and counter taken directly from key size + 16 bytes of entropy XOR the personalization string | Key and counter derived from entropy XOR the personalization string by the update function (§10.2.1.2), starting from an all-zero key and counter | `Test_DRBG_BacktrackingResistance_CAVP` |
| **Personalization string length** (Table 3) | Longer strings are XORed in cyclically | At most key size + 16 bytes; otherwise `NewReader` returns `ErrPersonalizationTooLong` | `Test_NewReader_BacktrackingResistance_PersonalizationLength` |
| **Generate output** (§10.2.1.5.1 steps 3–5) | For each block, the counter is incremented and encrypted; the output is truncated to the requested length | Same | `Test_DRBG_BacktrackingResistance_MatchesSpec` |
| **Additional input in generate** (§10.2.1.5.1 step 2) | `ReadWithAdditionalInput` reseeds from entropy, the personalization string, and the additional input | Mixed into the key and counter by the update function without new entropy; at most key size + 16 bytes (`ErrAdditionalInputTooLong`); empty input is treated as absent | `Test_DRBG_BacktrackingResistance_CAVP_AdditionalInput`, `Test_DRBG_BacktrackingResistance_AdditionalInputSemantics` |
| **Update after generate** (§10.2.1.5.1 step 6) | Not performed; the key changes only on reseed or key rotation | A new key and counter are derived after every request | `Test_DRBG_BacktrackingResistance_MatchesSpec`, `Test_DRBG_BacktrackingResistance_CAVP` |
| **Reseed** (§10.2.1.4.1) | A new state derived from entropy, the personalization string, and the additional input replaces the old state | Entropy XOR the additional input is mixed into the existing key and counter; the personalization string is not used. Applies to `Reseed`, `ReseedInterval`, `ReseedRequests`, fork detection, and key rotation | `Test_DRBG_BacktrackingResistance_CAVP_Reseed`, `Test_DRBG_BacktrackingResistance_ReseedPaths` |
| **Prediction resistance** (§9.3.1 step 7) | Reseeds before every request; additional input is ignored | Reseeds per §10.2.1.4.1 before every request, mixing in the request's additional input, which is then not used again for that request | `Test_DRBG_BacktrackingResistance_CAVP_PredictionResistance` |
| **Maximum request size** (Table 3: 2^19 bits) | 64 KB; larger requests return `ErrRequestTooLarge` | Same | `Test_DRBG_Read_MaxRequestSize` |
| **Reseed interval** (Table 3: 2^48 requests) | Optional through `ReseedRequests` (0 disables it), capped at 2^48 | Enforced; `ReseedRequests` defaults to 2^48 when not set | `Test_NewReader_BacktrackingResistance_ReseedRequests`, `TestConfig_WithReseedRequests_Clamping` |
| **Health testing** (§11.3) | Not performed | Known-answer tests of the instantiate, reseed, and generate functions (with and without prediction resistance, for each key size) run once per process before first use; a failure returns `ErrDRBGSelfTestFailed` for the life of the process | `Test_RunDRBGKAT_Success`, `Test_NewReader_BacktrackingResistance_DRBGSelfTestFailure` |

## Related FIPS 140-2 Features

| Feature | Behavior | Tests |
|---|---|---|
| **AES known-answer test** (FIPS 140-2 §4.9.1) | `WithSelfTests` runs an AES-256-CTR known-answer test (NIST SP 800-38A F.5.5) once per process; a failure returns `ErrSelfTestFailed` | `Test_RunSelfTests_Success` |
| **Continuous random number generator test** (FIPS 140-2 §4.9.2) | `WithContinuousHealthTest` compares every generated 128-bit block with the previous one on the same instance, withholds the first block after creation, and enters a permanent error state (`ErrHealthTestFailed`) on a match | `Test_DRBG_ContinuousHealthTest`, `Test_DRBG_ContinuousHealthTest_FirstBlockWithheld` |
| **Key zeroization** (FIPS 140-2 §4.7.6) | `WithZeroization` zeroes the stored old key and counter during automatic key rotation; the AES key schedule held by `cipher.Block` and seed buffers are not erased | `Test_AsyncRekey_WithZeroization_Enabled` |

## Other Behavior

| Behavior | Details |
|---|---|
| **Fork detection** | After a process fork, the child reseeds before producing output. The process ID is checked on every request, or every `ForkDetectionInterval` requests. |
| **Concurrency** | Instances are pooled per shard (`WithShards`); each instance is used by one goroutine at a time, and `Read` is safe for concurrent use. |
| **Dependencies** | Go standard library only (`crypto/aes`, `crypto/cipher`, `crypto/rand`). |
