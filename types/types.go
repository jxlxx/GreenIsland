package types

import (
	"math/rand"

	"github.com/google/uuid"
)

type User struct {
	ID   uuid.UUID
	Name string
}

type Value struct {
	Value   int `yaml:"value"`
	Jitter  int `yaml:"jitter"`
	Average int `yaml:"average_delta"`
}

func (v Value) CalcUpdate() int {
	if v.Jitter == 0 {
		return 0
	}
	delta := rand.Intn(v.Jitter) - v.Jitter/2
	shift := v.Average
	return delta + shift

}
