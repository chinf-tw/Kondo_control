package main

import (
	"fmt"
	"log"

	"github.com/jacobsa/go-serial/serial"
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

	// Write 4 bytes to the port.
	b := []byte{0b10000000, 0x3A, 0x4C}
	n, err := port.Write(b)
	if err != nil {
		log.Fatalf("port.Write: %v", err)
	}

	fmt.Println("Wrote", n, "bytes.")
	data := make([]byte, 256)
	i, err := port.Read(data)
	if err != nil {
		log.Fatalf("port.Read: %v", err)
	}
	fmt.Println("Read", data[:i])
}
