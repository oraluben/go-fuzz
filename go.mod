module github.com/oraluben/go-fuzz

go 1.16

require (
	github.com/elazarl/go-bindata-assetfs v1.0.1
	github.com/pingcap/parser v0.0.0-20210427084954-8e8ed7927bde
	github.com/pingcap/tidb v1.1.0-beta.0.20210430090150-27cacd8caf64
	github.com/pragmatwice/go-squirrel v0.0.0-20210430104750-5ab14fa4fe3b
	github.com/stephens2424/writerset v1.0.2
	golang.org/x/tools v0.1.0
)

replace github.com/pragmatwice/go-squirrel => ../go-squirrel
