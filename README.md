Original README of go-fuzz has been renamed to `README.go-fuzz.md`

# Ti-fuzz: Fuzz TiDB!

## Usage

### Prerequisites

```shell
go env -w GOPRIVATE=github.com/pragmatwice/go-squirrel
git config --global url.git@github.com:.insteadOf https://github.com/
# Clone modified TiDB that supports fuzzing.
git clone git@github.com:oraluben/tidb.git
```

### Build & run

1. build go-fuzz-build
2. build go-fuzz
3. `cd <tidb-root>/tidb-server/fuzz`
4. `go-fuzz-build -o tidb-fuzz.zip`
5. `go-fuzz -bin tidb-fuzz.zip`
