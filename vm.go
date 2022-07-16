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
	ErrInvalidWEND  = errors.New("unexpected instruction WEND: WHILE was not before")
	ErrInvalidNEXT  = errors.New("unexpected instruction NEXT: stack pointer at max offset")
	ErrInvalidPREV  = errors.New("unexpected instruction PREV: stack pointer at zero offset")
	ErrInvalidPUT   = errors.New("unexpected instruction PUT: there is no output writer")
	ErrInvalidPULL  = errors.New("unexpected instruction PULL: there is no input reader")
)

// NewDefault creates new vm instance with default parameters.
func NewDefault(program []Inst) *VM {
	return New(os.Stdin, os.Stdout, program)
}

// New creates new vm instance.
func New(in io.Reader, out io.Writer, program []Inst) *VM {
	return &VM{
		inp:  in,
		out:  out,
		prog: program,
		cpu:  make([]byte, 30000),
	}
}

// VM is brainfuck virtual machine.
type VM struct {
	prog []Inst
	pc   int

	inp io.Reader
	out io.Writer

	jmp []int

	i   int
	cpu []byte
}

// Run starts vm.
func (vm *VM) Run() error {
	return vm.RunContext(context.Background())
}

// RunContext starts vm with context.
func (vm *VM) RunContext(ctx context.Context) (err error) {
	if vm.pc != 0 || vm.jmp != nil {
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
	skip := len(vm.jmp) > 0 && vm.jmp[len(vm.jmp)-1] == -1
	if exit = vm.pc >= len(vm.prog); exit {
		if len(vm.jmp) > 0 {
			err = ErrInvalidWHILE
		}
		return
	}
	inst := vm.prog[vm.pc]
	if skip && inst != WEND && inst != WHILE {
		vm.pc++
		return false, nil
	}
	switch inst {
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
		jmp := vm.pc
		if vm.cpu[vm.i] == 0 {
			jmp = -1
		}
		vm.jmp = append(vm.jmp, jmp)
	case WEND:
		if len(vm.jmp) == 0 {
			err = ErrInvalidWEND
			break
		}
		if vm.cpu[vm.i] == 0 {
			vm.jmp = vm.jmp[:len(vm.jmp)-1]
			break
		}
		vm.pc = vm.jmp[len(vm.jmp)-1]
	default:
	}

	vm.pc++
	return
}
