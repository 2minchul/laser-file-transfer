package constants

import "time"

const (
	ReadPinNumber  = 18 // Change this to your GPIO PIN number
	WritePinNumber = 17
	Delay          = 30 * time.Millisecond
	StartPattern   = 0b10000110
	StartDelay1    = 90 * time.Millisecond
	StartDelay2    = 45 * time.Millisecond
)
