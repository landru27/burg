package main

type Worker struct {
	name string
	job  Job
}

type Job interface {
	goToWork(*Stockpile)
	performJob(*Stockpile, chan int)
}
