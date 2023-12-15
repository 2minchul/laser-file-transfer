package constants

import "time"

const (
	ReadPinNumber  = 18 // Change this to your GPIO PIN number
	WritePinNumber = 17
	Delay          = 1000 * time.Millisecond
	StartPattern   = 0b01000101
)
