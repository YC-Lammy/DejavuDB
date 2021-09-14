package parser

import (
	"github.com/xwb1989/sqlparser"
)

func walkExpr(expr sqlparser.Expr){
	switch expr:= expr.(type){
	//case *sqlparser.AliasedExpr:
	//case *sqlparser.AliasedTableExpr:

	case *sqlparser.AndExpr:

	case *sqlparser.BinaryExpr:

	case *sqlparser.CaseExpr:

	case *sqlparser.CollateExpr:

	case *sqlparser.ComparisonExpr:

	case *sqlparser.ConvertExpr:

	case *sqlparser.ConvertUsingExpr:

	case *sqlparser.ExistsExpr:

	case *sqlparser.FuncExpr:

	case *sqlparser.GroupConcatExpr:

	case *sqlparser.IntervalExpr:

	case *sqlparser.IsExpr:

	//case *sqlparser.JoinTableExpr:

	case *sqlparser.MatchExpr:

	case *sqlparser.NotExpr:

	case *sqlparser.OrExpr:

	case *sqlparser.ParenExpr:

	//case *sqlparser.ParenTableExpr:

	//case *sqlparser.SelectExpr:

	//case *sqlparser.SetExpr:

	//case *sqlparser.StarExpr:

	case *sqlparser.SubstrExpr:

	//case *sqlparser.TableExpr:

	case *sqlparser.UnaryExpr:

	//case *sqlparser.UpdateExpr:

	case *sqlparser.ValuesFuncExpr:

	}
}