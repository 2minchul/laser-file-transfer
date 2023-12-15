package main

import (
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"time"

	"github.com/flaviostutz/signalutils"
	"github.com/stianeikeland/go-rpio/v4"

	"laser/constants"
	"laser/protocol"
)

type StateChangeEvent struct {
	IsUpperRange bool
	Time         time.Time
}

func main() {
	err := rpio.Open()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer rpio.Close()

	pin := rpio.Pin(constants.ReadPinNumber)
	pin.Input()

	st, _ := signalutils.NewSchmittTrigger(30, 70, true)

	go func() {
		// Continuously read from the pin
		for {
			// Read pin state
			level := pin.Read()
			value := float64(level * 100)
			st.SetCurrentValue(value)
		}
	}()

	getCurrentState := func() bool {
		return st.IsUpperRange()
	}

	ch := make(chan *StateChangeEvent, 10)

	go func() {
		previous := getCurrentState()
		for {
			state := getCurrentState()
			if state != previous {
				ch <- &StateChangeEvent{IsUpperRange: state, Time: time.Now()}
				previous = state
			}
		}
	}()

	reader, writer := io.Pipe()
	go func() {
		previousTime := time.Now()
		var b byte
		var i int
		for {
			select {
			case event := <-ch:
				diff := event.Time.Sub(previousTime)
				previousTime = event.Time
				fmt.Println(diff, float64(diff/constants.Delay))
				if diff > time.Minute {
					fmt.Println("skip... no data in 1 min")
					i = 0
					b = 0
					continue
				}
				cnt := math.Round(float64(diff / constants.Delay))
				if cnt == 0 {
					continue
				}

				char := "0"
				if event.IsUpperRange {
					char = "1"
				}
				fmt.Println(strings.Repeat(char, int(cnt)))

				// event.IsUpperRange 가 cnt 만큼 반복된 bit 를 byte 로 만들어서 writer 에 쓴다.
				for j := 0; j < int(cnt); j++ {
					b = b << 1
					if event.IsUpperRange {
						b = b | 1
					}
					i++
					if i == 8 {
						writer.Write([]byte{b})
						i = 0
						b = 0
					}
				}

			}
		}
	}()

	for {
		fmt.Println("waiting ...")
		m, err := protocol.ReadFileMessage(reader)
		if err != nil {
			fmt.Println(err)
			continue
		}

		func() {
			f, err := os.Create(m.FileName)
			if err != nil {
				err = fmt.Errorf("failed to create file %s: %w", m.FileName, err)
				fmt.Println(err)
				return
			}
			defer f.Close()
			_, err = f.Write(m.Content)
			if err != nil {
				err = fmt.Errorf("failed to write to file %s: %w", m.FileName, err)
				fmt.Println(err)
				return
			}
			fmt.Printf("file %s saved!\n", m.FileName)
		}()
	}
}
