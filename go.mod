module github.com/oraluben/go-fuzz

go 1.16

require (
	github.com/elazarl/go-bindata-assetfs v1.0.1
	github.com/pingcap/parser v0.0.0-20210421190254-588138d35e55
	github.com/pingcap/tidb v2.0.11+incompatible
	github.com/stephens2424/writerset v1.0.2
	golang.org/x/tools v0.1.0
)

replace (
	github.com/oraluben/go-fuzz => ./
	github.com/pingcap/tidb v2.0.11+incompatible => github.com/oraluben/tidb v1.1.0-beta.0.20210428004842-82327ebaba06
)
