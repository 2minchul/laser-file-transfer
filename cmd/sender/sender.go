package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/stianeikeland/go-rpio/v4"

	"laser/constants"
	"laser/protocol"
)

func main() {
	filename := os.Args[1]

	msg, err := func() (*protocol.FileMessage, error) {
		f, err := os.Open(filename)
		if err != nil {
			err = fmt.Errorf("cannot open file: %w", err)
			return nil, err
		}
		defer f.Close()
		content, err := io.ReadAll(f)
		if err != nil {
			err = fmt.Errorf("cannot read file: %w", err)
			return nil, err
		}
		msg := &protocol.FileMessage{
			FileNameSize: len(filename),
			FileName:     filename,
			ContentSize:  uint64(len(content)),
			Content:      content,
		}

		return msg, nil
	}()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = rpio.Open()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer rpio.Close()
	pin := rpio.Pin(constants.WritePinNumber)
	pin.Output()

	reader, writer := io.Pipe()
	go func() {
		defer writer.Close()
		_, err := msg.WriteTo(writer)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	// 데이터 전송
	transmitData(reader, pin)
	println()
}

func transmitData(reader io.Reader, pin rpio.Pin) {
	pin.Toggle()
	pin.Toggle()
	pin.High()
	time.Sleep(constants.StartDelay1)
	pin.Low()
	time.Sleep(constants.StartDelay2)

	buffer := make([]byte, 1)
	for {
		n, err := reader.Read(buffer)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Read error:", err)
			}
			break
		}
		if n > 0 {
			sendByte(pin, buffer[0])
		}
	}
	pin.Toggle()
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
