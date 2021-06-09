package vm

func Compile(code []byte) []Inst {
	comp := make([]Inst, len(code))

	for i := range code {
		comp[i] = CompileOne(code[i])
	}

	return comp
}

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
		return END
	}
}
