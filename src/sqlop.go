package main

type Vdbe struct {
	opcode byte

	instruction []byte

	int_p  *int
	int8_p *int8
	str_p  *string
}
