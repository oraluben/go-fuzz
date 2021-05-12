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

### Dump & visualization

`go-fuzz ... -dumpcover` will generate `coverprofile`
`go tool cover -html=coverprofile`

Note: `go-fuzz` will not always generate valid coverage file for `go tool cover`, you might
need `sed -i"" -e '/0.0,1.1/d' coverprofile` (MacOS) or `sed -i '/0.0,1.1/d' coverprofile` (Linux) before generating
HTML coverage report.
