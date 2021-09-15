package main

import (
	"fmt"
	"log"

	"github.com/xwb1989/sqlparser"
)

func visit(node sqlparser.SQLNode) (bool, error) {
	fmt.Printf("%T \n", node)
	return true, nil
}

func main() {
	sql := `SELECT city, state_province, country_id
FROM locations
WHERE country_id IN('UK', 'CA');`
	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		// Do something with the err
		log.Fatal(err)
	}

	//sqlparser.Walk(visit, stmt)

	_, err = visit(stmt)
	if err != nil {
		panic(err)
	}
	node := stmt.(*sqlparser.Select)

	var v1 = ""
	var v2 = ""
	var v3 = ""
	var v4 = ""
	var v5 = ""
	var v6 = ""
	var v7 = ""
	var v8 = ""
	sqlparser.Walk(func(node sqlparser.SQLNode) (bool, error) {
		fmt.Printf(v1+"%T \n", node)
		v1 += " "
		return true, nil
	}, node.Comments)
	sqlparser.Walk(func(node sqlparser.SQLNode) (bool, error) {
		fmt.Printf(v2+"%T \n", node)
		v2 += " "
		return true, nil
	}, node.SelectExprs)
	sqlparser.Walk(func(node sqlparser.SQLNode) (bool, error) {
		fmt.Printf(v3+"%T \n", node)
		v3 += " "
		return true, nil
	}, node.From)
	sqlparser.Walk(func(node sqlparser.SQLNode) (bool, error) {
		fmt.Printf(v4+"%T \n", node)
		v4 += " "
		return true, nil
	}, node.Where)
	sqlparser.Walk(func(node sqlparser.SQLNode) (bool, error) {
		fmt.Printf(v5+"%T \n", node)
		v5 += " "
		return true, nil
	}, node.GroupBy)
	sqlparser.Walk(func(node sqlparser.SQLNode) (bool, error) {
		fmt.Printf(v6+"%T \n", node)
		v6 += " "
		return true, nil
	}, node.Having)
	sqlparser.Walk(func(node sqlparser.SQLNode) (bool, error) {
		fmt.Printf(v7+"%T \n", node)
		v7 += " "
		return true, nil
	}, node.OrderBy)
	sqlparser.Walk(func(node sqlparser.SQLNode) (bool, error) {
		fmt.Printf(v8+"%T \n", node)
		v8 += " "
		return true, nil
	}, node.Limit)

	// Otherwise do something with stmt

}
