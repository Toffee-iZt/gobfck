package vm

// Inst represents brainfuck instruction.
type Inst uint8

// instructions set.
const (
	NOP Inst = iota
	NEXT
	PREV
	INC
	DEC
	PUT
	PULL
	WHILE
	WEND
)

func (i *Inst) String() string {
	switch *i {
	case NEXT:
		return "NEXT"
	case PREV:
		return "PREV"
	case INC:
		return "INC"
	case DEC:
		return "DEC"
	case PUT:
		return "PUT"
	case PULL:
		return "PULL"
	case WHILE:
		return "WHILE"
	case WEND:
		return "WEND"
	default:
		return "NOP"
	}
}
