package main

import (
	"log"

	"github.com/wangle201210/zskc/internal/originalgame"
)

func main() {
	if err := originalgame.Run(); err != nil {
		log.Fatal(err)
	}
}
