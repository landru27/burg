package main

import (
	"math/rand"
	"time"
)

type Wheatfarmer struct {
	feedback chan int
}

func (wf *Wheatfarmer) goToWork(stockpile *Stockpile) {
	wf.feedback = make(chan int)

	randgen := rand.New(rand.NewSource(time.Now().UnixNano()))

	for {
		wf.performJob(stockpile, wf.feedback)

		time.Sleep(time.Duration(randgen.Intn(1500)+250) * time.Millisecond)
	}
}

func (wf *Wheatfarmer) performJob(stockpile *Stockpile, feedback chan int) {
	farmedwheat := Stockupdate{
		itemname: "wheat",
		itemqty:  80,
		result:   feedback,
	}
	stockpile.dropoff <- farmedwheat
	_ = <-feedback
}
