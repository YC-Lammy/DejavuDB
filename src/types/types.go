package types

type DType byte

const (
	String DType = iota
	Int
	Int8
	Int16
	Int32
	Int64
	Int128
	Uint
	Uint64
	Uint32
	Uint16
	Uint8
	Uint128
	Float
	Float32
	Float64
	Float128
	Decimal
	Decimal32
	Decimal64
	Decimal128
	Byte
	Byte_arr
	Bool
	Graph
	Table
	Json
	SmartContract
	Contract
	Money
	SmallMoney
	Time
	Date
	Datetime
	Smalldatetime

	Array_interface
	Null
	Map_string_interface
)

/* standard switch

case types.String:
case types.Int, types.Int64:
case types.Int32:
case types.Int16:
case types.Int8:
case types.Int128:
case types.Uint, types.Uint64:
case types.Uint32:
case types.Uint16:
case types.Uint8:
case types.Uint128:
case types.Decimal, types.Decimal64:
case types.Decimal32:
case types.Decimal128:
case types.Float, types.Float64:
case types.Float32:
case types.Float128:
case types.Byte:
case types.Byte_arr:
case types.Bool:
case types.Graph:
case types.Table:
case types.Json:
case types.SmartContract:
case types.Contract:
case types.Money:
case types.SmallMoney:
case types.Time:
case types.Date:
case types.Datetime:
case types.Smalldatetime:
case types.Null:

*/
