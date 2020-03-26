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

		time.Sleep(time.Duration(randgen.Intn(1500)+250) * time.Millisecond)
	}
}

func (m *Miller) performJob(stockpile *Stockpile, feedback chan int) {
	wheattogrind := Stockupdate{
		itemname: "wheat",
		itemqty:  12,
		result:   feedback,
	}
	stockpile.pickup <- wheattogrind
	pickedup := <-feedback

	if pickedup < 12 {
		fmt.Printf("miller could only find %d wheat!\n", pickedup)
		return
	}

	groundflour := Stockupdate{
		itemname: "flour",
		itemqty:  36,
		result:   feedback,
	}
	stockpile.dropoff <- groundflour
	_ = <-feedback
}
