package main

type Vdbe struct {
	opcode byte

	int_p  *int
	int8_p *int8
	str_p  *string
}
