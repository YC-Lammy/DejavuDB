package astwalker

import "github.com/xwb1989/sqlparser"

func Union(stmt *sqlparser.Union) (*Result, error) {
	var result = new(Result)
	switch stmt := stmt.Left.(type) {
	case *sqlparser.SELECT:
		Select(stmt)
	case *sqlparser.ParenSelect:
	case *sqlparser.Union:
		Union(stmt)
	}
	switch stmt := stmt.Right.(type) {
	case *sqlparser.SELECT:
		Select(stmt)
	case *sqlparser.ParenSelect:
	case *sqlparser.Union:
		Union(stmt)
	}

	switch stmt.Type {
	case "union":
	case "union all":
	case "uion distinct":
	}

	result, err := OrderBy(result, stmt.OrderBy)
	if err != nil {
		return nil, err
	}
	return result, nil
}
