package serial

import (
	"bytes"
	"fmt"
	"io"
	"kondocontrol/internal/convert"
	"kondocontrol/internal/eeprom"

	"github.com/pkg/errors"
)

type SubCommand uint8

const (
	readByteLength = 68
)
const (
	ScEEPROM      SubCommand = 0x00
	ScStretch     SubCommand = 0x01
	ScSpeed       SubCommand = 0x02
	ScCurrent     SubCommand = 0x03
	ScTemperature SubCommand = 0x04
)

// WriteEEPROM
func WriteEEPROM(id uint8, sc SubCommand, data []byte, port io.ReadWriteCloser) ([]byte, error) {
	var (
		cmd uint8 = 0b11000000 + id
	)
	if sc == ScEEPROM {
		// Confirm that this data is normal EEPROM data
		_, err := eeprom.Parse(data)
		if err != nil {
			return nil, errors.Wrap(err, "[WriteEEPROM]")
		}
	}
	b := []byte{cmd, uint8(sc)}
	b = append(b, data...)
	result, err := writeAndRead(port, b)
	if err != nil {
		return nil, errors.Wrap(err, "[WriteEEPROM]")
	}
	return result, nil
}

// ReadEEPROM
func ReadEEPROM(id uint8, sc SubCommand, port io.ReadWriteCloser) ([]byte, error) {
	var (
		cmd uint8 = 0b10100000 + id
	)
	b := []byte{cmd, uint8(sc)}
	fmt.Println(printHex(b))
	result, err := writeAndRead(port, b)
	if err != nil {
		return nil, errors.Wrap(err, "[ReadEEPROM]")
	}
	if len(result) < 2 {
		return nil, errors.New("The result length should not be smaller than 2")
	}
	// Confirm that this data is normal EEPROM data
	_, err = eeprom.Parse(result[2:])
	if err != nil {

		return nil, errors.Wrap(err, "[ReadEEPROM]")
	}
	return result[2:], nil
}

// SetPosition
func SetPosition(id uint8, target uint, port io.ReadWriteCloser) (uint, error) {
	position := convert.New(target)
	cmd := byte(0b10000000) + id
	b := []byte{cmd, position.PosH, position.PosL}
	result, err := writeAndRead(port, b)
	if err != nil {
		return 0, err
	}
	if len(result) != 3 {
		return 0, errors.New("Waring this return data is worng")
	}
	tchH := result[1]
	tchL := result[2]
	r := convert.Position{PosH: tchH, PosL: tchL}
	return r.PosToUint(), nil
}

// SetFree
func SetFree(id uint8, port io.ReadWriteCloser) (uint, error) {
	var cmd byte = 0b10000000 + id
	b := []byte{cmd, 0, 0}
	result, err := writeAndRead(port, b)
	if err != nil {
		return 0, err
	}
	tchH := result[1]
	tchL := result[2]
	r := convert.Position{PosH: tchH, PosL: tchL}
	return r.PosToUint(), nil
}

func writeAndRead(port io.ReadWriteCloser, b []byte) ([]byte, error) {
	writeN, err := port.Write(b)
	if err != nil {
		return []byte{}, errors.Wrap(err, "[WriteAndRead] port.Write")
	}
	if writeN != len(b) {
		return []byte{}, errors.New(
			"[WriteAndRead] prot.write data length is not equaly origin data length")
	}
	data := make([]byte, readByteLength)
	readN, err := port.Read(data)
	if err != nil {
		return []byte{}, errors.Wrap(err, "[WriteAndRead] port.Read")
	}
	if readN < len(b) {
		return []byte{}, errors.New(
			"[WriteAndRead] Read data should be bigger than length b")
	}
	if !bytes.Equal(data[:len(b)], b) {
		return []byte{}, errors.New(
			"[WriteAndRead] These should equal")
	}
	return data[writeN:readN], nil
}

func printHex(bs []byte) string {
	sum := ""
	for _, b := range bs {
		sum += fmt.Sprintf("%X", b) + " "
	}
	return sum
}
