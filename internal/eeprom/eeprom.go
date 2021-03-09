package eeprom

import (
	"fmt"

	"github.com/pkg/errors"
)

// EEPROM is KONDO servo motor eeprom
type EEPROM struct {
	StretchGain                  uint8
	Speed                        uint8
	Punch                        uint8
	DeadBand                     uint8
	Damping                      uint8
	SafeTimer                    uint8
	Flag                         Flag
	MaximumPulseLimit            uint16
	MinimumPulseLimit            uint16
	SignalSpeed                  SignalSpeed
	TemperatureLimit             uint8
	CurrentLimit                 uint8
	Response                     uint8
	UserOffset                   int8
	ID                           uint8
	CharacteristicChangeStretch1 uint8
	CharacteristicChangeStretch2 uint8
	CharacteristicChangeStretch3 uint8
	Address                      Address
}

// Flag is KONDO servo motor flag
type Flag struct {
	Reverse      bool
	Free         bool
	PWMINH       bool
	RotationMode bool
	SlaveMode    bool
}
type SignalSpeed uint32

const (
	Low  SignalSpeed = 115200
	Mid  SignalSpeed = 625000
	High SignalSpeed = 1250000
)

var (
	// ErrDataLength is when the data length is not 66
	ErrDataLength = errors.New("The data length is not 66")
	// ErrDataMismatch is when The data is mismatch
	ErrDataMismatch = errors.New("The data is mismatch")
)

// Parsing resolve bytes to EEPROM
func Parsing(bs []byte) (EEPROM, error) {
	if len(bs) != 64 {
		return EEPROM{}, ErrDataLength
	}
	var (
		result = EEPROM{}
		mark   = uint8(0)
		// recordAddress = make(map[string][]byte)
		address = Address{}
	)

	recordFunc := func(key *Interval, add uint8) {
		*key = NewInterval(mark, mark+add)
		mark += add
	}

	{ // Fixed as 0x5A
		beginning, err := sliceByteToUint8(bs[mark : mark+2])
		if err != nil {
			return EEPROM{}, errors.WithStack(err)
		}
		if beginning != 0x5A {
			return EEPROM{},
				errors.Wrapf(
					ErrDataMismatch,
					"Target is %X, but actual %X, origin: %v",
					0x5A, beginning, bs)
		}
		recordFunc(&address.Fixed, 2)
	}
	{ // Stretch gain
		// 2,4…254 2-step sequence
		// consisting only of even
		stretchGain, err := sliceByteToUint8(bs[mark : mark+2])
		if err != nil {
			return EEPROM{}, errors.WithStack(err)
		}
		if stretchGain%2 != 0 {
			return EEPROM{}, errors.WithStack(ErrDataMismatch)
		}
		result.StretchGain = stretchGain
		recordFunc(&address.StretchGain, 2)
	}
	{ // Speed
		// 1,2,3…127
		speed, err := sliceByteToUint8(bs[mark : mark+2])
		if err != nil {
			return EEPROM{}, errors.WithStack(err)
		}
		if speed > 127 {
			return EEPROM{}, errors.WithStack(ErrDataMismatch)
		}
		result.Speed = speed
		recordFunc(&address.Speed, 2)
	}
	{ // Punch
		// 0,1,2,3…10
		punch, err := sliceByteToUint8(bs[mark : mark+2])
		if err != nil {
			return EEPROM{}, errors.WithStack(err)
		}
		if punch > 10 {
			return EEPROM{}, errors.WithStack(ErrDataMismatch)
		}
		result.Punch = punch
		recordFunc(&address.Punch, 2)
	}
	{ // Dead band
		// 0,1,2,3,4,5
		deadBand, err := sliceByteToUint8(bs[mark : mark+2])
		if err != nil {
			return EEPROM{}, errors.WithStack(err)
		}
		if deadBand > 10 { //FIXME: that need to check
			return EEPROM{},
				errors.Wrapf(
					ErrDataMismatch,
					"Dead band cann't bigger then 5, but actual is %d, data is %v, origin: %v, result: %v",
					bs[mark:mark+2], deadBand, bs, result)
		}
		result.DeadBand = deadBand
		recordFunc(&address.DeadBand, 2)
	}
	{ // Damping
		// 1,2…255
		damping, err := sliceByteToUint8(bs[mark : mark+2])
		if err != nil {
			return EEPROM{}, errors.WithStack(err)
		}
		if damping == 0 {
			return EEPROM{}, errors.WithStack(ErrDataMismatch)
		}
		result.Damping = damping
		recordFunc(&address.Damping, 2)
	}
	{ // Safe timer
		// 10,11…255
		// (0x01-0xFF)
		safeTimer, err := sliceByteToUint8(bs[mark : mark+2])
		if err != nil {
			return EEPROM{}, errors.WithStack(err)
		}
		if safeTimer == 0 {
			return EEPROM{}, errors.WithStack(ErrDataMismatch)
		}
		result.SafeTimer = safeTimer
		recordFunc(&address.SafeTimer, 2)
	}
	{ // Flag
		flagDetail := bs[mark : mark+2]
		if flagDetail[0]&0xF0 != 0 || flagDetail[1]&0xF0 != 0 {
			return EEPROM{}, errors.WithStack(ErrDataMismatch)
		}
		if flagDetail[0]&0b00000110 != 0 {
			return EEPROM{}, errors.WithStack(ErrDataMismatch)
		}
		if (flagDetail[1] & 0b00000100 >> 2) != 1 {
			return EEPROM{}, errors.Wrapf(ErrDataMismatch, "flagDetail[1]&0b00000100 != 1,actual: %v, data: %v", flagDetail[1]&0b00000100, flagDetail)
		}
		flag := Flag{
			SlaveMode:    flagDetail[0]&0b00001000>>3 == 1,
			RotationMode: flagDetail[0]&0b00000001 == 1,
			PWMINH:       flagDetail[1]&0b00001000>>3 == 1,
			Free:         flagDetail[1]&0b00000010>>1 == 1,
			Reverse:      flagDetail[1]&0b00000001 == 1,
		}
		result.Flag = flag
		recordFunc(&address.Flag, 2)
	}
	{ // Maximum pulse limit
		maximumPulseLimit, err := sliceByteToUint16(bs[mark : mark+4])
		if err != nil {
			return EEPROM{}, errors.WithStack(err)
		}
		if maximumPulseLimit < 3500 || maximumPulseLimit > 11500 {
			return EEPROM{}, errors.WithStack(ErrDataMismatch)
		}
		result.MaximumPulseLimit = maximumPulseLimit
		recordFunc(&address.MaximumPulseLimit, 4)
	}
	{ // Minimum pulse limit
		minimumPulseLimit, err := sliceByteToUint16(bs[mark : mark+4])
		if err != nil {
			return EEPROM{}, errors.WithStack(err)
		}
		if minimumPulseLimit < 3500 || minimumPulseLimit > 11500 {
			return EEPROM{}, errors.WithStack(ErrDataMismatch)
		}
		result.MinimumPulseLimit = minimumPulseLimit
		recordFunc(&address.MinimumPulseLimit, 4)
	}
	mark += 2
	{ // Signal speed
		ss, err := sliceByteToUint8(bs[mark : mark+2])
		if err != nil {
			return EEPROM{}, errors.WithStack(err)
		}
		switch ss {
		case 0x00:
			result.SignalSpeed = High
		case 0x01:
			result.SignalSpeed = Mid
		case 0x10:
			result.SignalSpeed = Low
		}
		recordFunc(&address.SignalSpeed, 2)
	}
	{ // Temperature limit
		temperatureLimit, err := sliceByteToUint8(bs[mark : mark+2])
		if err != nil {
			return EEPROM{}, errors.WithStack(err)
		}
		if temperatureLimit < 1 || temperatureLimit > 127 {
			return EEPROM{}, errors.WithStack(ErrDataMismatch)
		}
		result.TemperatureLimit = temperatureLimit
		recordFunc(&address.TemperatureLimit, 2)
	}
	{ // Current limit
		currentLimit, err := sliceByteToUint8(bs[mark : mark+2])
		if err != nil {
			return EEPROM{}, errors.WithStack(err)
		}
		if currentLimit < 1 || currentLimit > 63 {
			return EEPROM{}, errors.WithStack(ErrDataMismatch)
		}
		result.CurrentLimit = currentLimit
		recordFunc(&address.CurrentLimit, 2)
	}
	for i := 0; i < 9; i++ {
		mark += 2
	}
	{ // Response
		response, err := sliceByteToUint8(bs[mark : mark+2])
		if err != nil {
			return EEPROM{}, errors.WithStack(err)
		}
		if response < 1 || response > 5 {
			return EEPROM{}, errors.WithStack(ErrDataMismatch)
		}
		result.Response = response
		recordFunc(&address.Response, 2)
	}
	{ // User offset
		fmt.Println(bs[mark : mark+2])
		userOffset, err := sliceByteToInt8(bs[mark : mark+2])
		if err != nil {
			return EEPROM{}, errors.WithStack(err)
		}
		result.UserOffset = userOffset
		recordFunc(&address.UserOffset, 2)
	}
	mark += 2
	{ // ID
		id, err := sliceByteToUint8(bs[mark : mark+2])
		if err != nil {
			return EEPROM{}, errors.WithStack(err)
		}
		if id > 31 {
			return EEPROM{}, errors.WithStack(ErrDataMismatch)
		}
		result.ID = id
		recordFunc(&address.ID, 2)
	}
	{ // Characteristic change stretch 1
		stretch1, err := sliceByteToUint8(bs[mark : mark+2])
		if err != nil {
			return EEPROM{}, errors.WithStack(err)
		}
		if stretch1 == 0 || stretch1%2 == 1 {
			return EEPROM{}, errors.Wrapf(ErrDataMismatch, "actual: %d", stretch1)
		}
		result.CharacteristicChangeStretch1 = stretch1
		recordFunc(&address.CharacteristicChangeStretch1, 2)
	}
	{ // Characteristic change stretch 2
		stretch2, err := sliceByteToUint8(bs[mark : mark+2])
		if err != nil {
			return EEPROM{}, errors.WithStack(err)
		}
		if stretch2 == 0 || stretch2%2 == 1 {
			return EEPROM{}, errors.WithStack(ErrDataMismatch)
		}
		result.CharacteristicChangeStretch2 = stretch2
		recordFunc(&address.CharacteristicChangeStretch2, 2)
	}
	{ // Characteristic change stretch 3
		stretch3, err := sliceByteToUint8(bs[mark : mark+2])
		if err != nil {
			return EEPROM{}, errors.WithStack(err)
		}
		if stretch3 == 0 || stretch3%2 == 1 {
			return EEPROM{}, errors.WithStack(ErrDataMismatch)
		}
		result.CharacteristicChangeStretch3 = stretch3
		recordFunc(&address.CharacteristicChangeStretch3, 2)
	}
	result.Address = address
	return result, nil
}
func sliceByteToUint8(bs []byte) (uint8, error) {
	if len(bs) != 2 {
		return 0, ErrDataLength
	}
	return uint8((bs[0] << 4) + bs[1]), nil
}
func sliceByteToInt8(bs []byte) (int8, error) {
	if len(bs) != 2 {
		return 0, ErrDataLength
	}
	value := (bs[0]<<4)&0b01110000 + bs[1]
	if (bs[0]<<4)&0b10000000 == 0 {
		return -int8(value), nil
	}
	return int8(value), nil
}
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
func print(target map[string][]byte) {
	for k, v := range target {
		fmt.Printf("*** %s ***\n", k)
		for i := v[0]; i < v[len(v)-1]; i++ {
			fmt.Print(i, " ")
		}
		println()
	}
}
