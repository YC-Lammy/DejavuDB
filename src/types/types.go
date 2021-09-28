package types

const (
	String        = 0x00
	Int           = 0x01
	Int8          = 0x02
	Int16         = 0x03
	Int32         = 0x04
	Int64         = 0x05
	Int128        = 0x06
	Uint          = 0x07
	Uint64        = 0x08
	Uint32        = 0x09
	Uint16        = 0x0A
	Uint8         = 0x0B
	Uint128       = 0x0C
	Float         = 0x0D
	Float32       = 0x0E
	Float64       = 0x0F
	Float128      = 0x10
	Decimal       = 0x11
	Decimal32     = 0x12
	Decimal64     = 0x13
	Decimal128    = 0x14
	Byte          = 0x15
	Byte_arr      = 0x16
	Bool          = 0x17
	Graph         = 0x18
	Table         = 0x19
	Json          = 0x1A
	SmartContract = 0x1B
	Contract      = 0x1C
	Money         = 0x1D
	SmallMoney    = 0x1E
	Time          = 0x1F
	Date          = 0x20
	Datetime      = 0x21
	Smalldatetime = 0x22

	Array_interface      = 0x30
	Null                 = 0x31
	Map_string_interface = 0x32
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
