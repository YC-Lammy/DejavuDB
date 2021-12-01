package astwalker

import "github.com/xwb1989/sqlparser"

func Select(stmt *sqlparser.Select) (*Table, error) {
	result := new(Result)
	var tables = map[string]*Table{}
	TableParents := map[string]map[string]*Table{
		"default": tables,
	}
	for _, tableExpr := range stmt.From {
		switch expr := tableExpr.(type) {
		case *sqlparser.AliasedTableExpr:
			expr.Expr
			var table *Table
			if expr.As.String() != "" {
				tables[expr.As.String()] = table
			}
		case *sqlparser.ParenTableExpr:
		case *sqlparser.JoinTableExpr:
		}
	}

	var rs = []*Result{}
	for _, selectExpr := range stmt.SelectExprs {
		switch expr := selectExpr.(type) {
		case *sqlparser.AliasedExpr: // AliasedExpr defines an aliased SELECT expression.
		case *sqlparser.StarExpr: // StarExpr defines a '*' or 'table.*' expression.
		case sqlparser.Nextval: // Nextval defines the NEXT VALUE expression.
			r, err := Expr(expr.Expr)
			if err != nil {
				return nil, err
			}
			rs = rs.append(rs, r)
		case *sqlparser.ColName:
		}
	}
	for _, v := range rs {
		switch v.Type {

		}
	}

	result.Type = Table

	return result, nil
}
