package table

type Table struct {
	Name       string
	Columns    map[string]*Column
	permission [3]uint8
}

type Column struct {
	Name  string
	Dtype byte
	Data  []interface{}
}
