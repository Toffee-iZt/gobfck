package vm

import (
	"errors"
	"fmt"
	"io"
)

// New creates new vm instance.
func New(in io.Reader, out io.Writer, instStream InstStream) *VM {
	return &VM{
		inp:        in,
		out:        out,
		instStream: instStream,
	}
}

// VM is brainfuck virtual machine.
type VM struct {
	err error

	inp io.Reader
	out io.Writer

	instStream InstStream

	i   int
	cpu [30000]byte
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
	var brc int
	var jmp []int
	var inst []Inst
	var iptr int

	for ins := END; vm.err == nil; iptr++ {
		if iptr > len(inst)-1 {
			ins = vm.instStream.Next()
			inst = append(inst, ins)
		} else {
			ins = inst[iptr]
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
			if vm.cpu[vm.i] != 0 && brc == len(jmp) {
				jmp = append(jmp, iptr)
			}
			brc++
		case WEND:
			if brc == 0 {
				vm.err = errors.New("unexpected instruction WEND: WHILE was not before")
				break
			}
			if brc == len(jmp) {
				iptr = jmp[len(jmp)-1] - 1
				jmp = jmp[:len(jmp)-1]
			}
			brc--

		default:
			break
		}
	}

	vm.instStream = nil

	return vm.err
}
