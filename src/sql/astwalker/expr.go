package astwalker

import (
	"errors"

	"github.com/xwb1989/sqlparser"
)

func Expr(expr sqlparser.Expr, table ...*Result) (*Result, error) {
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

		v1, ok1 := r1.Data.(bool)
		v2, ok2 := r2.Data.(bool)
		if ok1 && ok2 {

			result.Data = v && v1

		} else {
			t := sqlparser.NewTrackedBuffer(nil)
			expr.Format(t)
			return nil, errors.New(`syntax error near "` + t.String() + `": OR expected boolean`)
		}
		result.Type = Value
		result.Dtype = Bool
		return result, nil

	case *sqlparser.OrExpr:
		r1, err := Expr(expr.Left)
		if err != nil {
			return nil, err
		}
		r2, err := Expr(expr.Right)
		if err != nil {
			return nil, err
		}
		v1, ok1 := r1.Data.(bool)
		v2, ok2 := r2.Data.(bool)
		if ok1 && ok2 {

			result.Data = v || v1

		} else {
			t := sqlparser.NewTrackedBuffer(nil)
			expr.Format(t)
			return nil, errors.New(`syntax error near "` + t.String() + `": OR expected boolean`)
		}
		result.Type = Value
		result.Dtype = Bool
		return result, nil

	case *sqlparser.NotExpr:
		r1, err := Expr(expr.Expr)
		if err != nil {
			return nil, err
		}
		v1, ok := r1.Data.(bool)
		if ok {

			result.Data = !v1

		} else {
			t := sqlparser.NewTrackedBuffer(nil)
			expr.Format(t)
			return nil, errors.New(`syntax error near "` + t.String() + `": OR expected boolean`)
		}
		result.Type = Value
		result.Dtype = Bool
		return result, nil

	case *sqlparser.ParenExpr: // parenthesized boolean
		return Expr(expr.Expr)

	case *sqlparser.ComparisonExpr:
		r1, err := Expr(expr.Left)
		if err != nil {
			return nil, err
		}
		r2, err := Expr(expr.Right)
		if err != nil {
			return nil, err
		}
		r3, err := Expr(expr.Escape)
		if err != nil {
			return nil, err
		}

		switch r1.Type {
		case Column_Name:
			if r2.Type == r1.Type {
				switch expr.Operater {
				case "=":
				case "<":
				case ">":
				case "<=":
				case ">=":
				case "!=":
				case "<=>":
				case "in":
				case "not in":
				case "like":
				case "not like":
				case "regexp":
				case "not regexp":
				case "->":
				case "->>":
				}

			} else {
				switch expr.Operater {
				case "=":
				case "<":
				case ">":
				case "<=":
				case ">=":
				case "!=":
				case "<=>":
				case "in":
				case "not in":
				case "like":
				case "not like":
				case "regexp":
				case "not regexp":
				case "->":
				case "->>":
				}
			}

		default:
			switch expr.Operater {
			case "=":
			case "<":
			case ">":
			case "<=":
			case ">=":
			case "!=":
			case "<=>":
			case "in":
			case "not in":
			case "like":
			case "not like":
			case "regexp":
			case "not regexp":
			case "->":
			case "->>":
			}
		}

	case *sqlparser.RangeCond: // RangeCond represents a BETWEEN or a NOT BETWEEN expression.
		r1, err := Expr(expr.Left)
		if err != nil {
			return nil, err
		}
		r2, err := Expr(expr.From)
		if err != nil {
			return nil, err
		}
		r3, err := Expr(expr.To)
		if err != nil {
			return nil, err
		}
		switch expr.Operater {
		case "between":
		case "not between":
		}

	case *sqlparser.IsExpr: // IsExpr represents an IS ... or an IS NOT ... expression.

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
	default:
	}
	return result, nil
}
