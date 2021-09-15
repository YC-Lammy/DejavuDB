package astwalker

const (
	Value       = 'V'
	Column      = 'C'
	Column_Name = 'N'
	Data_type   = 'D'
	Table       = 'T' // *Table, pointer to the actual table
	Table_Name  = 'H'
)

const (
	Bool   = 0x00
	Null   = 0x01
	String = 0x02
	Int    = 0x03
)

type Result struct {
	Type  byte
	Dtype byte
	Data  interface{}
	As    string
}
