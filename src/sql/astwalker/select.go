package astwalker

import "github.com/xwb1989/sqlparser"

func Select(stmt *sqlparser.Select) (*Result, error) {
	result := new(Result)
	var tables = map[string]*Result{}
	switch expr := stmt.(type) {
	case *sqlparser.AliasedTableExpr:
	case *sqlparser.ParenTableExpr:
	case *sqlparser.JoinTableExpr:
	}
	var rs = []*Result{}
	for _, v := range stmt.SelectExprs {
		switch expr := stmt.SelectExpr.(type) {
		case *sqlparser.AliasedExpr: // AliasedExpr defines an aliased SELECT expression.
		case *sqlparser.StarExpr: // StarExpr defines a '*' or 'table.*' expression.
		case sqlparser.Nextval: // Nextval defines the NEXT VALUE expression.
			r, err := Expr(expr.Expr)
			if err != nil {
				return nil, err
			}
			rs = rs.append(rs, r)
		}
	}
	for _, v := range rs {
		switch v.Type {

		}
	}

	result.Type = Table

	return result, nil
}
