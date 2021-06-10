package vm

func Compile(code []byte, dst []Inst) []Inst {
	for i := range code {
		dst = append(dst, CompileOne(code[i]))
	}

	return dst
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
