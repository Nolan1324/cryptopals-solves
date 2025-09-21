package main

import (
	"cryptopals/internal/apps/timeattack"
	"cryptopals/internal/randx"
	"time"
)

func main() {
	server := timeattack.NewServer(randx.RandBytes(16), time.Millisecond*2)
	server.Run("localhost:9000")
}
