package main

import (
	"os"

	"github.com/Toffee-iZt/gobfck/vm"
)

func main() {
	if len(os.Args) < 2 {
		println("specify the path to the file")
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		println(err.Error())
	}

	prog := vm.Compile(data)

	v := vm.New(os.Stdin, os.Stdout, prog)
	err = v.Run()
	if err != nil {
		println(err.Error())
	}
}
