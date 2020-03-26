package main

const minimumMillisecondsForJob = 250
const rangeMillisecondsForJob = 1500

type Worker struct {
	name string
	job  Job
}

type Job interface {
	goToWork(*Stockpile)
	performJob(*Stockpile, chan int)
}
