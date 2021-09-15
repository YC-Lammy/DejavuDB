package astwalker

import "github.com/xwb1989/sqlparser"

func Limit(result *Result, stmt *sqlparser.Limit) (*Result, error) {
	return result, nil
}
