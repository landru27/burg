package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Printf("starting\n")

	burgermeister := &Burgermeister{}
	burgermeister.initializeBurg()

	go burgermeister.updateStockpile()

	for {
		burgermeister.recruitWorkers()
		burgermeister.feedWorkers()

		burgermeister.listWorkers()
		burgermeister.showStockpile()

		time.Sleep(2000 * time.Millisecond)
	}

	fmt.Printf("stopping\n")
}
