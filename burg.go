package main

import (
	"fmt"
	"time"
)

const dailyCycleMilliseconds = 2000

func main() {
	fmt.Printf("starting\n")
	fmt.Printf("\n")

	// instantiate someone to run the town, and set up the town
	burgermeister := &Burgermeister{}
	burgermeister.initializeBurg()

	// the burgermeister handles additions to and removals from and queries about the stockpile;
	// kick off a routine to continuously receive and process those actions
	go burgermeister.manageStockpile()

	// the burgermeister's daily cycle : recruit, feed, and report
	for {
		burgermeister.recruitWorkers()
		burgermeister.feedWorkers()

		burgermeister.listWorkers()
		burgermeister.showStockpile()

		time.Sleep(dailyCycleMilliseconds * time.Millisecond)
	}

	fmt.Printf("\n")
	fmt.Printf("stopping\n")
}
