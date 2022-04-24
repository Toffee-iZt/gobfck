package vm

import (
	"errors"
	"fmt"
	"io"
)

// errors.
var (
	ErrInvalidInst  = errors.New("invalid instruction")
	ErrInvalidWHILE = errors.New("invalid instruction WHILE: no WEND instruction found")
	ErrInvalidWEND  = errors.New("unexpected operator WEND: WHILE was not before")
	ErrInvalidNEXT  = errors.New("invalid instruction NEXT: stack pointer at max offset")
	ErrInvalidPREV  = errors.New("invalid instruction PREV: stack pointer at zero offset")
)

// New creates new vm instance.
func New(in io.Reader, out io.Writer, program []Inst) *VM {
	return &VM{
		inp:  in,
		out:  out,
		prog: program,
	}
}

// VM is brainfuck virtual machine.
type VM struct {
	err error

	prog []Inst
	pc   int

	inp io.Reader
	out io.Writer

	jmp []int
	skp []bool

	i   int
	cpu [30000]byte
}

// Run starts vm.
func (vm *VM) Run() error {
	if vm.prog == nil {
		vm.err = errors.New("vm: no instructions for exec")
		return vm.err
	}
	for vm.err == nil {
		if vm.do() {
			break
		}
	}
	println()
	return vm.err
}

func (vm *VM) getch() (byte, error) {
	var c [1]byte
	_, err := vm.inp.Read(c[:])
	return c[0], err
}

func (vm *VM) print(c byte) {
	fmt.Fprintf(vm.out, "%c", c)
}

func (vm *VM) do() bool {
	skip := len(vm.skp) > 0 && vm.skp[len(vm.skp)-1]
	if vm.pc >= len(vm.prog) {
		if skip {
			vm.err = ErrInvalidWHILE
		}
		return true
	}
	inst := vm.prog[vm.pc]
	if skip && inst != WEND {
		vm.pc++
		return false
	}
	switch inst {
	case NEXT:
		if vm.i < len(vm.cpu)-1 {
			vm.i++
			break
		}
		vm.err = ErrInvalidNEXT
	case PREV:
		if vm.i > 0 {
			vm.i--
			break
		}
		vm.err = ErrInvalidPREV
	case INC:
		vm.cpu[vm.i]++
	case DEC:
		vm.cpu[vm.i]--
	case PUT:
		vm.print(vm.cpu[vm.i])
	case PULL:
		b, err := vm.getch()
		if err != nil {
			vm.err = err
		}
		vm.cpu[vm.i] = b
	case WHILE:
		vm.jmp = append(vm.jmp, vm.pc)
		vm.skp = append(vm.skp, vm.cpu[vm.i] == 0)
	case WEND:
		if len(vm.jmp) == 0 {
			vm.err = ErrInvalidWEND
			break
		}
		if vm.cpu[vm.i] == 0 {
			vm.jmp = vm.jmp[:len(vm.jmp)-1]
			vm.skp = vm.skp[:len(vm.skp)-1]
			break
		}
		vm.pc = vm.jmp[len(vm.jmp)-1]
	default:
		vm.err = ErrInvalidInst
	}

	vm.pc++
	return false
}
