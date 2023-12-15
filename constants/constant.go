package constants

import "time"

const (
	ReadPinNumber  = 18 // Change this to your GPIO PIN number
	WritePinNumber = 17
	Delay          = 100 * time.Millisecond
	StartPattern   = 0b01000101
)
