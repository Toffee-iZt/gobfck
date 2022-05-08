package gobfck

// Compile compiles code into vm instructions.
func Compile(code []byte) []Inst {
	program := make([]Inst, 0, len(code))
	for i := range code {
		program = append(program, CompileOne(code[i]))
	}
	return program
}

// CompileOne converts byte into instruction.
func CompileOne(b byte) Inst {
	switch b {
	case '>':
		return NEXT
	case '<':
		return PREV
	case '+':
		return INC
	case '-':
		return DEC
	case '.':
		return PUT
	case ',':
		return PULL
	case '[':
		return WHILE
	case ']':
		return WEND
	default:
		return NOP
	}
}
