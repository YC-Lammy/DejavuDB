package astwalker

type ValueType byte

const (
	IntType ValueType = iota
)

type Result struct {
	Type         byte
	Dtype        byte
	Data         interface{}
	As           string
	RowsAffected int
}

type Column struct {
	Type  ValueType
	Value []interface{}
}

type Value struct {
	Type  ValueType
	Value interface{}
}
