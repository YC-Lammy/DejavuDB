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
