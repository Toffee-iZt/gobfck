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
		run([]byte(hw))
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		println(err.Error())
	}
	run(data)
}

func run(data []byte) {
	v := gobfck.New(os.Stdin, os.Stdout, gobfck.Compile(data))
	err := v.Run()
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}
