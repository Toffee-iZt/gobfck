package gobfck

// Compile compiles code into vm instructions.
func Compile(code []byte) (Bytecode, error) {
	prog := make(Bytecode, 0, len(code))
	var stack []uint16
	for i := range code {
		c := OpcodeFor(code[i])
		var op uint16
		switch c {
		case NOP:
			continue
		case WHILE:
			stack = append(stack, uint16(len(prog)))
		case WEND:
			if len(stack) < 0 {
				return nil, ErrInvalidWEND
			}
			op = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			prog.SetOperand(int(op), uint16(len(prog)))
		}
		prog = prog.Append(c, op)
	}
	if len(stack) > 0 {
		return nil, ErrInvalidWHILE
	}
	return prog, nil
}

// OpcodeFor converts byte into opcode.
func OpcodeFor(b byte) Opcode {
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
