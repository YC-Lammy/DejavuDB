package astwalker

import "github.com/xwb1989/sqlparser"

type Table struct {
	name      string
	subTables map[string]*Table
	indexHint *sqlparser.IndexHints
}

func (t *Table) Select(colnames ...string) *Table {}

func (t *Table) Filter(colname, val string) *Table {}

func (t *Table) Where()

func TableExpr(expr sqlparser.TableExpr) (*Table, error) {
	switch expr := expr.(type) {
	case *sqlparser.AliasedTableExpr:

		table, err := SimpleTableExpr(expr.Expr)

		if err != nil {
			return nil, err
		}

		if expr.As.String() != "" {
			table.name = expr.As.String()
		}
		table.indexHint = expr.Hints

		return table, nil

	case *sqlparser.ParenTableExpr:
	case *sqlparser.JoinTableExpr:
	}
	return result, nil
}

func SimpleTableExpr(expr sqlparser.SimpleTableExpr) (*Table, error) {
	switch expr := expr.(type) {
	case sqlparser.TableName:
	case *sqlparser.Subquery:
		switch expr.Select.(type) {
		case *sqlparser.Select:
		case *sqlparser.Union:
		case *sqlparser.ParenSelect:
		}
	}
}
