package eeprom

// Kondo for uint8
func sliceByteToUint8(bs []byte) (uint8, error) {
	if len(bs) != 2 {
		return 0, ErrDataLength
	}
	return uint8((bs[0] << 4) + bs[1]), nil
}
func uint8ToSliceByte(u uint8) []byte {
	high := u & 0xF0 >> 4
	low := u & 0x0F
	return []byte{high, low}
}

// Kondo for int8
func sliceByteToInt8(bs []byte) (int8, error) {
	if len(bs) != 2 {
		return 0, ErrDataLength
	}
	value := (bs[0]<<4)&0b01110000 + bs[1]
	if (bs[0]<<4)&0b10000000 == 0 {
		return int8(value), nil
	}
	return -int8(value), nil
}
func int8ToSliceByte(i int8) []byte {
	var (
		high    byte
		low     byte
		signBit byte = 0
	)
	if i < 0 {
		i--
		i = ^i
		signBit = 1 << 3
	}
	high = (byte(i)&0xF0)>>4 + signBit
	low = byte(i) & 0x0F
	return []byte{high, low}
}

// Kondo for uint16
func sliceByteToUint16(bs []byte) (uint16, error) {
	if len(bs) != 4 {
		return 0, ErrDataLength
	}
	zero := uint16(bs[0]) << (4 * 3)
	one := uint16(bs[1]) << (4 * 2)
	two := uint16(bs[2]) << (4 * 1)
	three := uint16(bs[3]) << (4 * 0)
	return zero + one + two + three, nil
}
func uint16ToSliceByte(u uint16) []byte {
	zero := uint8(u>>(4*3)) & 0x0F
	one := uint8(u>>(4*2)) & 0x0F
	two := uint8(u>>(4*1)) & 0x0F
	three := uint8(u>>(4*0)) & 0x0F
	return []byte{zero, one, two, three}
}

// Kondo Flag
func sliceByteToFlag(b []byte) Flag {
	return Flag{
		SlaveMode:    b[0]&0b00001000>>3 == 1,
		RotationMode: b[0]&0b00000001 == 1,
		PWMINH:       b[1]&0b00001000>>3 == 1,
		Free:         b[1]&0b00000010>>1 == 1,
		Reverse:      b[1]&0b00000001 == 1,
	}
}
func flagToSliceByte(flag Flag) []byte {
	var (
		slaveMode    uint8
		rotationMode uint8
		pwminh       uint8
		free         uint8
		reverse      uint8
	)
	if flag.SlaveMode {
		slaveMode = 0b00001000
	}
	if flag.RotationMode {
		rotationMode = 0b00000001
	}
	if flag.PWMINH {
		pwminh = 0b00001000
	}
	if flag.Free {
		free = 0b00000010
	}
	if flag.Reverse {
		reverse = 0b00000001
	}
	return []byte{
		slaveMode + rotationMode,
		pwminh + free + reverse + 0b00000100,
	}
}

// SignalSpeed
// func sliceByteToSignalSpeed(b []byte) SignalSpeed {
// 	b[0] >> 3
// }
func signalSpeedToSliceByte(s SignalSpeed) []byte {
	return []byte{byte(s), 0}
}
