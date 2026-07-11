package main

import (
	"log"

	"github.com/wangle201210/zskc/internal/game"
)

func main() {
	if err := game.Run(); err != nil {
		log.Fatal(err)
	}
}
