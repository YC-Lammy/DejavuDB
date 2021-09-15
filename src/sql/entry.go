package sql

import (
	"io"
	"strings"

	"./astwalker"

	"github.com/xwb1989/sqlparser"
)

type result_table [][]string

func Process_sql(text string) (result_table, error) {
	r := strings.NewReader(text)

	tokens := sqlparser.NewTokenizer(r)
	for {
		stmt, err := sqlparser.ParseNext(tokens)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		// Do something with stmt or err.
		switch stmt := stmt.(type) { // top level statments

		case *sqlparser.Select:
			astwalker.Select(stmt)
		case *sqlparser.ParenSelect:
		case *sqlparser.Union:
		case *sqlparser.Stream:

		case *sqlparser.Insert:
		case *sqlparser.Update:
		case *sqlparser.Delete:

		case *sqlparser.Set:

		case *sqlparser.DBDDL: // create or drop database
		case *sqlparser.DDL: // create, alter, drop, rename or truncate

		case *sqlparser.Show:
		case *sqlparser.Use:
		case *sqlparser.Begin:
		case *sqlparser.Commit:
		case *sqlparser.Rollback:
		case *sqlparser.OtherRead:
		case *sqlparser.OtherAdmin:
		}
	}

	return nil, nil

}
