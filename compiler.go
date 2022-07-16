package gobfck

import (
	"io"
	"os"
)

// Compile compiles code into vm instructions.
func Compile(code []byte) []Inst {
	program := make([]Inst, 0, len(code))
	for i := range code {
		c := CompileOne(code[i])
		if c == NOP {
			continue
		}
		program = append(program, c)
	}
	return program
}

// CompileFile compiles file with code.
func CompileFile(path string) ([]Inst, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return CompileReader(f)
}

// CompileReader func.
func CompileReader(r io.Reader) ([]Inst, error) {
	prog := make([]Inst, 0, 128)
	for {
		var buf [64]byte
		n, err := r.Read(buf[:])
		if err != nil {
			return nil, err
		}
		for i := 0; i < n; i++ {
			c := CompileOne(buf[i])
			if c == NOP {
				continue
			}
			prog = append(prog, c)
		}
		if n < len(buf) {
			break
		}
	}
	return prog, nil
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
