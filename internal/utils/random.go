package utils

import (
	"math/rand"
	"time"
)

func GenerateRoomPin() string {
	rand.Seed(time.Now().UnixNano())
	pin := ""
	for i := 0; i < 4; i++ {
		pin += string('0' + rune(rand.Intn(10)))
	}
	return pin
}

