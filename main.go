package main

import (
	"io"
	"os"

	"github.com/Toffee-iZt/gobfck/vm"
)

func main() {
	if len(os.Args) == 1 {
		//run interpreter.
		return
	}

	if len(os.Args) == 2 {
		runFile(os.Args[1])
		return
	}

	switch os.Args[1] {
	case "run":
		runFile(os.Args[2])
	}
}

func runFile(path string) {
	insts := compileFile(path)
	v := vm.New(os.Stdin, os.Stdout, vm.NewInstStream(insts))

	v.Run()
}

func compileFile(path string) []vm.Inst {
	f, err := os.OpenFile(path, os.O_RDONLY, 0777)
	if err != nil {
		panic(err)
	}
	stat, err := f.Stat()
	if err != nil {
		panic(err)
	}

	insts := make([]vm.Inst, 0, stat.Size())

	var buf [256]byte
	for {
		n, err := f.Read(buf[:])
		if n == 0 {
			if err == io.EOF {
				return insts
			}
			panic(err)
		}
		insts = vm.Compile(buf[:n], insts)
	}
}
