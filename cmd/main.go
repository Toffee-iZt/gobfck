package main

import (
	"os"

	"github.com/Toffee-iZt/gobfck"
)

const hw = "++++++++++[>+++++++>++++++++++>+++<<<-]>++.>+.+++++++..+++.>++.<<+++++++++++++++.>.+++.------.--------.>+."

func main() {
	if len(os.Args) < 2 {
		println("file path is not specified")
		println("starting 'Hello World'\n")
		run(gobfck.Compile([]byte(hw)))
	}

	prog, err := gobfck.CompileFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	run(prog)
}

func run(prog []gobfck.Inst) {
	v := gobfck.NewDefault(prog)
	err := v.Run()
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}
