package vm

import (
	"errors"
	"fmt"
	"io"
)

// New creates new vm instance.
func New(in io.Reader, out io.Writer, instStream InstStream) *VM {
	return &VM{
		done:       make(chan struct{}),
		inp:        in,
		out:        out,
		instStream: instStream,
	}
}

// VM is brainfuck virtual machine.
type VM struct {
	err  error
	done chan struct{}

	instStream InstStream

	inp io.Reader
	out io.Writer

	jmp  []int
	inst []Inst
	iptr int
	brc  int

	i   int
	cpu [30000]byte
}

func (vm *VM) Done() <-chan struct{} {
	return vm.done
}

func (vm *VM) Error() error {
	return vm.err
}

// Run starts vm.
func (vm *VM) Run() error {
	if vm.err != nil {
		return vm.err
	}
	if vm.instStream == nil {
		vm.err = errors.New("vm: no instructions for exec")
		return vm.err
	}
	return vm.run()
}

func (vm *VM) run() error {
	for vm.err == nil {
		vm.do()
	}

	vm.instStream = nil

	return vm.err
}

func (vm *VM) do() {
	var ins Inst
	if vm.iptr > len(vm.inst)-1 {
		ins = vm.instStream.Next()
		vm.inst = append(vm.inst, ins)
	} else {
		ins = vm.inst[vm.iptr]
	}

	switch ins {
	case NEXT:
		if vm.i == len(vm.cpu)-1 {
			vm.err = errors.New("invalid instruction NEXT: stack pointer has max offset")
			break
		}
		vm.i++
	case PREV:
		if vm.i == 0 {
			vm.err = errors.New("invalid instruction PREV: stack pointer at 0 offset")
			break
		}
		vm.i--
	case INC:
		vm.cpu[vm.i]++
	case DEC:
		vm.cpu[vm.i]--
	case PUT:
		//_, vm.err = vm.out.Write([]byte{vm.cpu[vm.i]})
		fmt.Fprintf(vm.out, "%c ", vm.cpu[vm.i]+32)
	case PULL:
		_, vm.err = vm.inp.Read(vm.cpu[vm.i : vm.i+1])
	case WHILE:
		if vm.cpu[vm.i] != 0 && vm.brc == len(vm.jmp) {
			vm.jmp = append(vm.jmp, vm.iptr)
		}
		vm.brc++
	case WEND:
		if vm.brc == 0 {
			vm.err = errors.New("unexpected instruction WEND: WHILE was not before")
			break
		}
		if vm.brc == len(vm.jmp) {
			vm.iptr = vm.jmp[len(vm.jmp)-1] - 1
			vm.jmp = vm.jmp[:len(vm.jmp)-1]
		}
		vm.brc--
	default:
		break
	}

	vm.iptr++
}
