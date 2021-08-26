package main

import (
	"flag"
	"fmt"
	"log"

	"kondocontrol/internal/khr_3hv"

	"github.com/jacobsa/go-serial/serial"
)

func main() {
	var (
		lp = flag.String("left-port", "", "left port")
		rp = flag.String("right-port", "", "right port")
	)
	flag.Parse()
	if *lp == "" || *rp == "" {
		log.Fatalf("left and right port should not be empty, (lp: %s,rp: %s)", *lp, *rp)
	}
	// Set up leftOptions.
	leftOptions := serial.OpenOptions{
		PortName:          *lp,
		BaudRate:          1250000,
		DataBits:          8,
		StopBits:          1,
		MinimumReadSize:   3,
		ParityMode:        serial.PARITY_EVEN,
		RTSCTSFlowControl: false,
	}
	rightOptions := leftOptions
	rightOptions.PortName = *rp
	// Open the port.
	rightPort, err := serial.Open(rightOptions)
	if err != nil {
		log.Fatalf("rightPort.Open: %v", err)
	}
	leftPort, err := serial.Open(leftOptions)
	if err != nil {
		log.Fatalf("leftPort.Open: %v", err)
	}

	// Make sure to close it later.
	defer rightPort.Close()
	defer leftPort.Close()

	speedValue := uint8(50)

	// init robot
	robot, err := khr_3hv.DefaultRobotNum(leftPort, rightPort)
	if err != nil {
		log.Fatal(err)
	}
	for _, r := range robot {
		da, err := r.SetSpeed(speedValue)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("WriteEEPROM: ", da)
	}
}
