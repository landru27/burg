package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Baker struct {
	feedback chan int
}

func (m *Baker) goToWork(stockpile *Stockpile) {
	m.feedback = make(chan int)

	randgen := rand.New(rand.NewSource(time.Now().UnixNano()))

	for {
		m.performJob(stockpile, m.feedback)

		time.Sleep(time.Duration(randgen.Intn(1500)+250) * time.Millisecond)
	}
}

func (m *Baker) performJob(stockpile *Stockpile, feedback chan int) {
	flourforbread := Stockupdate{
		itemname: "flour",
		itemqty:  48,
		result:   feedback,
	}
	stockpile.pickup <- flourforbread
	pickedup := <-feedback

	if pickedup < 48 {
		fmt.Printf("baker could only find %d flour!\n", pickedup)
		return
	}

	bakedbread := Stockupdate{
		itemname: "bread",
		itemqty:  12,
		result:   feedback,
	}
	stockpile.dropoff <- bakedbread
	_ = <-feedback
}
