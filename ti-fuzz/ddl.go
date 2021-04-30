package ti_fuzz

import (
	"database/sql"

	ctx "github.com/pragmatwice/go-squirrel/instantiator"
)

var hardCodedDDL = "CREATE TABLE a;"

func Initialize(conn *sql.DB) {
	conn.Exec(hardCodedDDL)
}

func CurrentInstantiateContext() *ctx.InstantiateContext {
	return nil
}
