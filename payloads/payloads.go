package payloads

import (
	"encoding/json"
)

func Bytes(data interface{}) []byte {
	b, _ := json.Marshal(data)
	return b
}
