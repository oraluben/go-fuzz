package ti_fuzz

import "github.com/pingcap/parser/model"

const Seed = "CREATE TABLE t1 (i1 INTEGER, c1 CHAR);" +
	"INSERT INTO t1 VALUES (1, 'a'), (2, 'b'), (3, 'c');" +
	"SELECT 1 FROM t1 WHERE i1 = 1;"

var Scheme = map[model.CIStr][]model.CIStr{
	model.NewCIStr("t1"): {
		model.NewCIStr("i1"), model.NewCIStr("c1"),
	},
}

var Libs = []string{
	"select * from c where c.a = c.b",
	"select x, y from a order by z limit 1",
	"select x from a where x < y order by x",
	"select * from a where y != 1 and z > 100 group by x",
	"select * from (select * from a) b join (select * from c) d",
	"select count(a) from t where a = 1 and b + a > 100",
	"select t1.a, t1.b, t2.c from t1 inner join t2 on t1.a = t2.b + t2.c",
}
