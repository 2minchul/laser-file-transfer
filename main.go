package main

import (
	"fmt"
	"os"
	"time"

	"github.com/stianeikeland/go-rpio/v4"

	"laser/constants"
)

func main() {
	err := rpio.Open()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer rpio.Close()
	pin := rpio.Pin(constants.WritePinNumber)
	pin.Output()
	sendByte(pin, byte(constants.StartPattern))
}

func sendByte(pin rpio.Pin, b byte) {
	for i := 0; i < 8; i++ {
		bit := (b >> (7 - i)) & 1
		if bit == 1 {
			fmt.Print("1")
			pin.High()
		} else {
			fmt.Print("0")
			pin.Low()
		}
		time.Sleep(constants.Delay)
	}
}
