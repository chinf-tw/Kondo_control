package convert

// Position is a Kondo servo motor position value
type Position struct {
	PosH   uint8
	PosL   uint8
	Origin uint
}

// PosToUint is PosH and PosL convert to uint and update self Origin
func (p *Position) PosToUint() uint {
	high := uint(p.PosH)
	low := uint(p.PosL)
	p.Origin = high<<7 + low
	return p.Origin
}

// New does new a Position struct
func New(position uint) Position {
	p := Position{Origin: position}
	p.PosH, p.PosL = uintToPos(position)
	return p
}

func uintToPos(position uint) (uint8, uint8) {
	posH := uint8((position >> 7) & 0b01111111)
	posL := uint8(position & 0x01111111)
	return posH, posL
}
