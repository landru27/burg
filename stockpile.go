package main

type Stockpile struct {
	stock   map[string]int
	dropoff chan Stockupdate
	pickup  chan Stockupdate
	query   chan Stockquery
}

type Stockupdate struct {
	itemname string
	itemqty  int
	result   chan int
}

type Stockquery struct {
	itemname string
	result   chan int
}

var orderedStockpile = [...]string{
	"wheat",
	"flour",
	"bread",
}
