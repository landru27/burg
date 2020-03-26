package main

type Stockpile struct {
	stock   map[string]int
	dropoff chan Stockupdate
	pickup  chan Stockupdate
	query   chan string
}

type Stockupdate struct {
	itemname string
	itemqty  int
	result   chan int
}
