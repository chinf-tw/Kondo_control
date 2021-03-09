package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"kondocontrol/internal/convert"
	"kondocontrol/internal/eeprom"
	"log"
	"time"

	"github.com/jacobsa/go-serial/serial"
)

type subCommand uint8

const (
	readByteLength = 68
)
const (
	scEEPROM      subCommand = 0x00
	scStretch     subCommand = 0x01
	scSpeed       subCommand = 0x02
	scCurrent     subCommand = 0x03
	scTemperature subCommand = 0x04
)

func main() {
	// Set up options.
	options := serial.OpenOptions{
		PortName:          "COM7",
		BaudRate:          1250000,
		DataBits:          8,
		StopBits:          1,
		MinimumReadSize:   3,
		ParityMode:        serial.PARITY_EVEN,
		RTSCTSFlowControl: false,
	}

	// Open the port.
	port, err := serial.Open(options)
	if err != nil {
		log.Fatalf("serial.Open: %v", err)
	}

	// Make sure to close it later.
	defer port.Close()

	{ // test readEEPROM
		// r := readEEPROM(0, scEEPROM, port)
		// fmt.Println(printHex(r), len(r))
	}

	{ // test setPosition and Free mode
		// trySetPosition(port)
		// tryFree(port)
		// speedValue := uint8(50)
		// for i := 0; i <= 10; i++ {
		// 	writeEEPROM(uint8(i), scSpeed, []byte{speedValue}, port)
		// }
	}

	{ // test data file
		// testingFilePath := "./Ignore/data"
		// dat, err := ioutil.ReadFile(testingFilePath)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// fmt.Println(printHex(dat))
	}
	{ // convert data address to json
		testingFilePath := "./Ignore/data"
		dat, err := ioutil.ReadFile(testingFilePath)
		if err != nil {
			log.Fatal(err)
		}
		ee, err := eeprom.Parsing(dat)
		if err != nil {
			log.Fatalf("%+v\n", err)
		}
		j, err := json.Marshal(ee.Address)
		if err != nil {
			log.Fatalf("json.Marshal: %v", err)
		}
		fmt.Println(string(j))
	}

}
func writeEEPROM(id uint8, sc subCommand, data []byte, port io.ReadWriteCloser) {
	var (
		cmd uint8 = 0b11000000 + id
	)

	b := []byte{cmd, uint8(sc)}
	b = append(b, data...)
	result := writeAndRead(port, b)
	fmt.Println(result)
}
func readEEPROM(id uint8, sc subCommand, port io.ReadWriteCloser) []byte {
	var (
		cmd uint8 = 0b10100000 + id
	)
	b := []byte{cmd, uint8(sc)}
	fmt.Println(printHex(b))
	result := writeAndRead(port, b)
	e, err := eeprom.Parsing(result[2:])
	if err != nil {
		log.Fatalf("eeprom.Parsing: %+v", err)
		return nil
	}
	fmt.Printf("%+v\n", e)

	// save data
	// testingFilePath := "./Ignore/testing"
	// {
	// 	f, err := os.Create(testingFilePath)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	n, err := f.Write(result[2:])
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	println(n)
	// 	f.Close()
	// }
	return result
}
func trySetPosition(port io.ReadWriteCloser) {
	const (
		targetEnd   uint  = 8000
		targetStart uint  = 7000
		count       uint  = 10
		step              = (targetEnd - targetStart) / count
		stepTime          = time.Second / 2
		idStart     uint8 = 0
		idEnd       uint8 = 10
	)

	for index := uint(0); index < count; index++ {
		target := uint(7500)
		// step target
		// target = targetStart + index*step
		if index%2 == 0 {
			target = targetStart
		} else {
			target = targetEnd
		}
		// fmt.Printf("% 08b, % 08b\n", result.PosH, result.PosL)
		for i := idStart; i <= idEnd; i++ {
			setPosition(i, target, port)
		}
		// fmt.Println("Read", data[:i])
		time.Sleep(stepTime)
	}
}
func setPosition(id uint8, target uint, port io.ReadWriteCloser) {
	position := convert.New(target)
	cmd := byte(0b10000000) + id
	b := []byte{cmd, position.PosH, position.PosL}
	result := writeAndRead(port, b)
	tchH := result[1]
	tchL := result[2]
	r := convert.Position{PosH: tchH, PosL: tchL}
	fmt.Printf("ID: %d, Target: %d, Current: %d\n", id, target, r.PosToUint())
}
func tryFree(port io.ReadWriteCloser) {
	for i := 0; i <= 10; i++ {
		id := uint8(i)
		var cmd byte = 0b10000000 + id
		b := []byte{cmd, 0, 0}
		result := writeAndRead(port, b)
		tchH := result[1]
		tchL := result[2]
		r := convert.Position{PosH: tchH, PosL: tchL}
		fmt.Printf("ID: %d, Current: %d\n", id, r.PosToUint())
		time.Sleep(time.Second)
	}
}
func writeAndRead(port io.ReadWriteCloser, b []byte) []byte {
	writeN, err := port.Write(b)
	if err != nil {
		log.Fatalf("port.Write: %v", err)
	}
	if writeN != len(b) {
		log.Fatal("[self] prot.write is not equaly origin data")
	}
	// fmt.Println("Wrote", n, "bytes.")
	data := make([]byte, readByteLength)
	readN, err := port.Read(data)
	if err != nil {
		log.Fatalf("port.Read: %v", err)
	}
	if readN > len(b) {
		if !bytes.Equal(data[:len(b)], b) {
			println("[self] Dose'nt Equal!!")
		}
	}

	return data[writeN:readN]
}
func printHex(bs []byte) string {
	sum := ""
	for _, b := range bs {
		sum += fmt.Sprintf("%X", b) + " "
	}
	return sum
}
