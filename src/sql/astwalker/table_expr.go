package astwalker

import "github.com/xwb1989/sqlparser"

func TableExpr(expr *sqlparser.TableExpr) (*Result, error) {
	result = new(Result)
	switch expr := expr.(type) {
	case *sqlparser.AliasedTableExpr:
	case *sqlparser.ParenTableExpr:
	case *sqlparser.JoinTableExpr:
	}
	return result, nil
}
