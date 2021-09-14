package main

import (
	"fmt"
	"log"
	"github.com/xwb1989/sqlparser"
)

func visit(node sqlparser.SQLNode) (bool, error){
	fmt.Printf("%T \n", node)
	return true, nil
}

func main() {
	sql := "SELECT a FROM table1 WHERE a = 'abc'"
stmt, err := sqlparser.Parse(sql)
if err != nil {
	// Do something with the err
	log.Fatal(err)
}

sqlparser.Walk(visit, stmt)
fmt.Println(stmt)

// Otherwise do something with stmt
switch stmt := stmt.(type) {
	
case  *sqlparser.AliasedExpr:


case *sqlparser.Select:
	_ = stmt

case *sqlparser.Stream:

case *sqlparser.Insert:

case *sqlparser.Update:

case *sqlparser.Delete:

case *sqlparser.DDL:

case *sqlparser.Begin:

case *sqlparser.Commit:

case *sqlparser.Rollback:

case *sqlparser.Set:

case *sqlparser.Show:

case *sqlparser.Use:

case *sqlparser.Union:
}
}
