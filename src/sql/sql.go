package sql

/* syntax check
// data key check
// seperate process columns
// register clauses
// register functions
// interprete to native nosql and golang func
// register task flow nodes
// execute nodes in order
// todo: optimization by monitering each node

type node struct { // use a node struct to clearify task order
	name string
	next *node
}
*/

const (
	sql_str8    = 0x00 // byte
	sql_str     = 0x01 // []string
	sql_int     = 0x02 // []int32
	sql_int64   = 0x03 // []int64
	sql_float32 = 0x04 // []float32
	sql_float64 = 0x05 // []float64
	sql_bool    = 0x06 // []bool
)

//
var readOnlyPermission = [3]int8{7, 5, 5}
var publicPermission = [3]int8{7, 7, 7}
var groupPermission = [3]int8{7, 7, 0}
var ownerPermission = [3]int8{7, 0, 0}

//

var SQLDB = map[string]*sql_database{}
var defaultSQLDB *sql_database

type sql_database struct {
	name           string // name of database
	schemas        map[string]*schema
	permission     [3]int8
	owner          int
	group          int
	default_schema *schema // pointer to the schema, default :dbo
}

type schema struct {
	name string

	tables        map[string]*table
	procedures    map[string]*procedure
	default_table *table

	permission [3]int8
	owner      int
	group      int
}

type views struct {
}

type procedure struct {
	hash   string
	opcode *Query
}

func init() {
	is := schema{}
	dbo := schema{permission: publicPermission}

	db1 := sql_database{
		name:           "Database_1",
		schemas:        map[string]*schema{"dbo": &dbo, "information_schema": &is},
		permission:     publicPermission,
		owner:          1000,
		group:          1000,
		default_schema: &dbo,
	}

	SQLDB["Database_1"] = &db1
	defaultSQLDB = &db1
}
