package vm

// Inst represents brainfuck instruction.
type Inst uint8

// instructions set.
const (
	NEXT Inst = iota
	PREV
	INC
	DEC
	PUT
	PULL
	WHILE
	WEND
	END
)

// InstStream is instructions stream for vm.
type InstStream interface {
	Next() Inst
}

type instArray struct {
	array []Inst
	ptr   int
}

func (b *instArray) Next() Inst {
	if b.ptr > len(b.array)-1 {
		return END
	}
	i := b.array[b.ptr]
	b.ptr++
	return i
}

// NewInstStream creates new instruction stream
// from compiled instructions array.
func NewInstStream(i []Inst) InstStream {
	return &instArray{array: i}
}
