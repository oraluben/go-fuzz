package ti_fuzz

import (
	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/model"
	"github.com/pragmatwice/go-squirrel/instantiator"
)

var Seed string
var Libs []string
var Scheme *instantiator.TableInfoContext

func GetScheme(p *parser.Parser, sql string) (*instantiator.TableInfoContext, error) {
	stmts, _, err := p.Parse(sql, "", "")
	if err != nil {
		return nil, err
	}

	ctx := instantiator.NewTableInfoContext(map[model.CIStr][]model.CIStr{})

	for _, v := range stmts {
		if ct, ok := v.(*ast.CreateTableStmt); ok {
			ctx.Build(ct)
		}
	}

	return ctx, nil
}
