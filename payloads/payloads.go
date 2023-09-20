package payloads

import "encoding/json"

type WorldTick struct {
	Quarter int
	Day     int
	Hour    int
}

func Bytes(data interface{}) []byte {
	b, _ := json.Marshal(data)
	return b
}
