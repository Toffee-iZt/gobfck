package gobfck

import (
	"context"
	"errors"
	"io"
	"os"
)

// errors.
var (
	ErrInvalidWHILE = errors.New("invalid instruction WHILE: no WEND instruction found")
	ErrInvalidWEND  = errors.New("invalid instruction WEND: WHILE was not before")
	ErrInvalidNEXT  = errors.New("unexpected instruction NEXT: stack pointer at max offset")
	ErrInvalidPREV  = errors.New("unexpected instruction PREV: stack pointer at zero offset")
	ErrInvalidPUT   = errors.New("unexpected instruction PUT: there is no output writer")
	ErrInvalidPULL  = errors.New("unexpected instruction PULL: there is no input reader")
)

// NewDefault creates new vm instance with default parameters.
func NewDefault(prog Bytecode) *VM {
	return New(os.Stdin, os.Stdout, prog)
}

// New creates new vm instance.
func New(in io.Reader, out io.Writer, prog Bytecode) *VM {
	return &VM{
		inp:  in,
		out:  out,
		prog: prog,
		cpu:  make([]byte, 30000),
	}
}

// VM is brainfuck virtual machine.
type VM struct {
	prog Bytecode
	pc   int

	inp io.Reader
	out io.Writer

	i   int
	cpu []byte
}

// Run starts vm.
func (vm *VM) Run() error {
	return vm.RunContext(context.Background())
}

// RunContext starts vm with context.
func (vm *VM) RunContext(ctx context.Context) (err error) {
	if vm.pc != 0 || vm.i != 0 {
		return errors.New("vm already holds state")
	}
	for exit := false; !exit && err == nil; {
		select {
		case <-ctx.Done():
			err = ctx.Err()
		default:
			exit, err = vm.do()
		}
	}
	vm.print('\n')
	return
}

func (vm *VM) getch() (byte, error) {
	if vm.inp == nil {
		return 0, ErrInvalidPULL
	}
	var c [1]byte
	_, err := vm.inp.Read(c[:])
	return c[0], err
}

func (vm *VM) print(c byte) error {
	if vm.out == nil {
		return ErrInvalidPUT
	}
	_, err := vm.out.Write([]byte{c})
	return err
}

func (vm *VM) do() (exit bool, err error) {
	if exit = vm.pc >= len(vm.prog); exit {
		return
	}
	opcode, operand := vm.prog.Read(vm.pc)
	switch opcode {
	case NEXT:
		if vm.i >= len(vm.cpu)-1 {
			err = ErrInvalidNEXT
			break
		}
		vm.i++
	case PREV:
		if vm.i <= 0 {
			err = ErrInvalidPREV
			break
		}
		vm.i--
	case INC:
		vm.cpu[vm.i]++
	case DEC:
		vm.cpu[vm.i]--
	case PUT:
		err = vm.print(vm.cpu[vm.i])
	case PULL:
		b, e := vm.getch()
		if e != nil {
			err = e
			break
		}
		vm.cpu[vm.i] = b
	case WHILE:
		if vm.cpu[vm.i] == 0 {
			vm.pc = int(operand)
		}
	case WEND:
		if vm.cpu[vm.i] != 0 {
			vm.pc = int(operand)
		}
	default:
	}

	vm.pc = vm.prog.Next(vm.pc)

	return
}
