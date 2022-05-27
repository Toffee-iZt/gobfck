package gobfck

import (
	"context"
	"errors"
	"io"
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

	i   int
	cpu [30000]byte
}

// Run starts vm.
func (vm *VM) Run() error {
	return vm.RunContext(context.Background())
}

// RunContext starts vm with context.
func (vm *VM) RunContext(ctx context.Context) error {
	if vm.prog == nil {
		vm.err = errors.New("vm: no instructions for exec")
		return vm.err
	}
	for vm.err == nil {
		compl := false
		select {
		case <-ctx.Done():
			vm.err = ctx.Err()
		default:
			compl = vm.do()
		}
		if compl {
			break
		}
	}
	vm.out.Write([]byte{'\n'})
	return vm.err
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

func (vm *VM) do() bool {
	skip := len(vm.jmp) > 0 && vm.jmp[len(vm.jmp)-1] == -1
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
		if vm.i >= len(vm.cpu)-1 {
			vm.err = ErrInvalidNEXT
			break
		}
		vm.i++
	case PREV:
		if vm.i <= 0 {
			vm.err = ErrInvalidPREV
			break
		}
		vm.i--
	case INC:
		vm.cpu[vm.i]++
	case DEC:
		vm.cpu[vm.i]--
	case PUT:
		vm.err = vm.print(vm.cpu[vm.i])
	case PULL:
		b, err := vm.getch()
		if err != nil {
			vm.err = err
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
			vm.err = ErrInvalidWEND
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
	return false
}
