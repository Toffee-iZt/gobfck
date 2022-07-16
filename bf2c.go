package gobfck

import (
	"io"
	"strings"
)

// CompileToC converts brainfuck code to equivalent C code.
func CompileToC(input io.Reader, w io.Writer) error {
	io.WriteString(w, "// This source was automatically generated with\n")
	io.WriteString(w, "// gobfck brainfuck compiler.\n\n")

	io.WriteString(w, "#include <stdio.h>\n\n")
	io.WriteString(w, "int main() {\n\tchar a[30000], *ptr = a;\n")

	for t := 0; ; {
		var b [1]byte
		_, err := input.Read(b[:])
		if err != nil {
			if err == io.EOF {
				io.WriteString(w, "\treturn 0;\n}")
				return nil
			}
			return err
		}
		var c string
		switch b[0] {
		case '>':
			c = "ptr++;"
		case '<':
			c = "ptr--;"
		case '+':
			c = "++*ptr;"
		case '-':
			c = "--*ptr;"
		case '[':
			c = "while (*ptr) {"
		case ']':
			c = "}"
			t--
		case '.':
			c = "putchar(*ptr);"
		case ',':
			c = "*ptr = getchar();"
		default:
			continue
		}

		io.WriteString(w, strings.Repeat("\t", t+1)+c+"\n")
		if b[0] == '[' {
			t++
		}
	}
}
