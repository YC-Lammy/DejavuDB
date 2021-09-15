package astwalker

import (
	"errors"

	"github.com/xwb1989/sqlparser"
)

func Expr(expr *sqlparser.Expr) (*Result, error) {
	var result = new(Result)
	switch expr := expr.(type) {
	case *sqlparser.AndExpr:
		r1, err := Expr(expr.Left)
		if err != nil {
			return nil, err
		}
		r2, err := Expr(expr.Right)
		if err != nil {
			return nil, err
		}
		if v, ok := r1.data.(bool); ok {
			if v1, ok := r2.data.(bool); ok {
				if v && v1 {
					result.data = true
				} else {
					result.data = false
				}
			} else {
				return nil, errors.New("")
			}
		}
	case *sqlparser.OrExpr:
	case *sqlparser.NotExpr:
	case *sqlparser.ParenExpr:
	case *sqlparser.ComparisonExpr:
	case *sqlparser.RangeCond:
	case *sqlparser.IsExpr:
	case *sqlparser.ExistsExpr:
	case *sqlparser.SQLVal:
	case *sqlparser.NullVal:
	case sqlparser.BoolVal:
	case *sqlparser.ColName:
	case sqlparser.ValTuple:
	case *sqlparser.Subquery:
	case sqlparser.ListArg:
	case *sqlparser.BinaryExpr:
	case *sqlparser.UnaryExpr:
	case *sqlparser.IntervalExpr:
	case *sqlparser.CollateExpr:
	case *sqlparser.FuncExpr:
	case *sqlparser.CaseExpr:
	case *sqlparser.ValuesFuncExpr:
	case *sqlparser.ConvertExpr:
	case *sqlparser.SubstrExpr:
	case *sqlparser.ConvertUsingExpr:
	case *sqlparser.MatchExpr:
	case *sqlparser.GroupConcatExpr:
	case *sqlparser.Default:
	}
	return result, nil
}
