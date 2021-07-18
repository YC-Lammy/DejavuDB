package main

type Vdbe struct {
	opcode byte

	instruction []byte

	table_p   *table
	table1_p  *table
	column_p  *column
	column1_p *column
	row_p     *row
	row1_p    *row

	int_p     *int
	int8_p    *int8
	str_p     *string
	str8_p    *byte
	float32_p *float32
	float64_p *float64
	bool_p    *bool
	blob_p    *BLOB

	int1_p     *int
	int81_p    *int8
	str1_p     *string
	str81_p    *byte
	float321_p *float32
	float641_p *float64
	bool1_p    *bool
	blob1_p    *BLOB

	next *Vdbe
	prev *Vdbe
}

type Query struct {
	first   *Vdbe
	last    *Vdbe
	current *Vdbe

	all []*Vdbe

	program_counter int

	virtual_tables map[string]*table
}
