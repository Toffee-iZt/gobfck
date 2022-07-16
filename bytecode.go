package gobfck

// Bytecode type.
type Bytecode []byte

// Append func.
func (bc Bytecode) Append(opcode Opcode, operand uint16) Bytecode {
	bc = append(bc, byte(opcode))
	if opcode&OpMask != 0 {
		bc = append(bc, byte(operand>>8), byte(operand&0xff))
	}
	return bc
}

// SetOperand func.
func (bc Bytecode) SetOperand(ptr int, operand uint16) bool {
	if Opcode(bc[ptr])&OpMask == 0 {
		return false
	}
	bc[ptr+1], bc[ptr+2] = byte(operand>>8), byte(operand&0xff)
	return true
}

// Next func.
func (bc Bytecode) Next(ptr int) int {
	if Opcode(bc[ptr])&OpMask != 0 {
		ptr += 2
	}
	return ptr + 1
}

// Read func.
func (bc Bytecode) Read(ptr int) (Opcode, uint16) {
	opcode := Opcode(bc[ptr])
	var operand uint16
	if opcode&OpMask != 0 {
		operand = uint16(bc[ptr+1])<<8 | uint16(bc[ptr+2])
	}
	return opcode, operand
}

// Opcode represents brainfuck operator.
type Opcode uint8

// instructions set.
const (
	NOP Opcode = iota
	NEXT
	PREV
	INC
	DEC
	PUT
	PULL
	WHILE  = iota | OpMask
	WEND   = iota | OpMask
	OpMask = Opcode(1 << 7)
)

func (i *Opcode) String() string {
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
