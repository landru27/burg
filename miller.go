package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Miller struct {
	feedback chan int
}

func (m *Miller) goToWork(stockpile *Stockpile) {
	m.feedback = make(chan int)

	randgen := rand.New(rand.NewSource(time.Now().UnixNano()))

	for {
		m.performJob(stockpile, m.feedback)

		time.Sleep(time.Duration(randgen.Intn(rangeMillisecondsForJob)+minimumMillisecondsForJob) * time.Millisecond)
	}
}

func (m *Miller) performJob(stockpile *Stockpile, feedback chan int) {
	wheattogrind := Stockupdate{
		itemname: "wheat",
		itemqty:  wheatForFlour,
		result:   feedback,
	}
	stockpile.pickup <- wheattogrind
	pickedup := <-feedback

	if pickedup < wheatForFlour {
		fmt.Printf("miller could only find %d wheat!\n", pickedup)
		return
	}

	groundflour := Stockupdate{
		itemname: "flour",
		itemqty:  flourFromWheat,
		result:   feedback,
	}
	stockpile.dropoff <- groundflour
	_ = <-feedback
}
