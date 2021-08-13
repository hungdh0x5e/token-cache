# Token Cache

Cache token into memory to reusable. Token are consistent when used concurrently by multiple goroutines.

[![Go Reference](https://pkg.go.dev/badge/github.com/hungdh0x5e/token-cache.svg)](https://pkg.go.dev/github.com/hungdh0x5e/token-cache)

### Usage

Step 1: Create custom token getter, implement *TokenGetter* interface

Step 2: Init token cache, add to Client

See [cache_test](cache_test.go) or [pkg.go.dev](https://pkg.go.dev/github.com/hungdh0x5e/token-cache) for further
documentation and examples.
