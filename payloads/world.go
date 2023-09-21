package payloads

import (
	"time"
)

type WorldTick struct {
	Quarter int
	Day     int
	Hour    int
	EGT     int
	ERT     time.Duration
}
