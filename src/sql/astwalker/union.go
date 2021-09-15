package astwalker

import "github.com/xwb1989/sqlparser"

func Union(stmt *sqlparser.Union) (*Result, error) {
	var result = new(Result)
	var err error

	var r1 *Result
	switch stmt := stmt.Left.(type) {
	case *sqlparser.SELECT:
		r1, err = Select(stmt)
	case *sqlparser.ParenSelect:
		r1, err = Select(stmt)
	case *sqlparser.Union:
		r1, err = Union(stmt)
	}

	var r2 *Result
	switch stmt := stmt.Right.(type) {
	case *sqlparser.SELECT:
		r2, err = Select(stmt)
	case *sqlparser.ParenSelect:
		r2, err = Select(stmt)
	case *sqlparser.Union:
		r2, err = Union(stmt)
	}

	switch stmt.Type {
	case "union":
	case "union all":
	case "uion distinct":
	}

	result, err = OrderBy(result, stmt.OrderBy)
	if err != nil {
		return nil, err
	}
	return result, nil
}
