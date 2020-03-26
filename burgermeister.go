package main

import (
	"fmt"
)

// set values that determine key determinations the burgermeister makes
const startingBread = 1000
const minimumBreadToRecruitBreadWorkers = startingBread * .95
const minimumBreadToRecruitGeneralWorkers = startingBread * 2

const breadForAMeal = 4

// Burgermeister models the person overseeing the wellbeing and operation of a town
type Burgermeister struct {
	stockpile        *Stockpile
	workers          []*Worker
	recruitmentRules []*RecruitmentRule
}

// RecruitmentRule models a determination for when to recruit whom into the town
type RecruitmentRule struct {
	ready  func(*Stockpile) bool
	worker *Worker
}

// initializeBurg establishes the initial structure and conditions for the town
func (bm *Burgermeister) initializeBurg() {
	// the town has a stockpile of resources which citizens contribute to and draw from
	bm.stockpile = &Stockpile{}

	// all we start with in the stockpile is enough bread to get us going
	bm.stockpile.stock = make(map[string]int, 0)
	bm.stockpile.stock["bread"] = startingBread

	// these are the channels citizens will use to contribute to and draw from the stockpile;
	// they are the stockpile's input and output conduits
	bm.stockpile.dropoff = make(chan Stockupdate)
	bm.stockpile.pickup = make(chan Stockupdate)
	bm.stockpile.query = make(chan Stockquery)

	// the burgermeister begins with no citizens :-(
	bm.workers = make([]*Worker, 0, 0)

	// the burgermeister has a set of rules for when to recruit new citizens
	bm.recruitmentRules = make([]*RecruitmentRule, 0, 0)
	answer := make(chan int)
	bm.recruitmentRules = append(bm.recruitmentRules, &RecruitmentRule{
		ready: func(stockpile *Stockpile) bool {
			stockpile.query <- Stockquery{"wheat", answer}
			amountWheat := <-answer
			if amountWheat < 10 {
				return true
			}

			return false
		},
		worker: &Worker{
			name: "wheat farmer",
			job:  &Wheatfarmer{},
		},
	})
	bm.recruitmentRules = append(bm.recruitmentRules, &RecruitmentRule{
		ready: func(stockpile *Stockpile) bool {
			stockpile.query <- Stockquery{"flour", answer}
			amountFlour := <-answer
			if amountFlour < 10 {
				return true
			}

			return false
		},
		worker: &Worker{
			name: "miller",
			job:  &Miller{},
		},
	})
	bm.recruitmentRules = append(bm.recruitmentRules, &RecruitmentRule{
		ready: func(stockpile *Stockpile) bool {
			stockpile.query <- Stockquery{"bread", answer}
			amountBread := <-answer
			if amountBread < 1000 {
				return true
			}

			return false
		},
		worker: &Worker{
			name: "baker",
			job:  &Baker{},
		},
	})
}

// recruitWorkers had its own global logic, and applies the logic encoded by the set of recruitmentRules
func (bm *Burgermeister) recruitWorkers() {
	if bm.stockpile.stock["bread"] < minimumBreadToRecruitBreadWorkers {
		fmt.Printf(">>>>  not enough bread to recruit new bread workers !\n")
		return
	}

	for _, v := range bm.recruitmentRules {
		if v.ready(bm.stockpile) {
			bm.recruitWorker(v.worker)
		}
	}
}

// the actual recruiting
func (bm *Burgermeister) recruitWorker(worker *Worker) {
	fmt.Printf("Burgermeister.recruitWorker : %s\n", worker.name)

	bm.workers = append(bm.workers, worker)
	go worker.job.goToWork(bm.stockpile)
}

// eat, drink, and keep working!
func (bm *Burgermeister) feedWorkers() {
	for _, v := range bm.workers {
		bm.feedWorker(v)
	}
}

// the actual feeding
func (bm *Burgermeister) feedWorker(worker *Worker) {
	feedback := make(chan int)

	eatsomebread := Stockupdate{
		itemname: "bread",
		itemqty:  breadForAMeal,
		result:   feedback,
	}
	bm.stockpile.pickup <- eatsomebread
	eaten := <-feedback

	if eaten < 1 {
		fmt.Printf(">>>>  WORKERS ARE STARVING !!!\n")
	} else if eaten < breadForAMeal {
		fmt.Printf(">>>>  workers are hungry !\n")
	}
}

// the burgermeister handles additions to and removals from the stockpile;
// this is the function by which we receive and process those updates;
// this function is designed to run continually, as a go routine
func (bm *Burgermeister) manageStockpile() {
	var added Stockupdate
	var taken Stockupdate
	var inquiry Stockquery

	for {
		select {
		case added = <-bm.stockpile.dropoff:
			bm.stockpile.stock[added.itemname] += added.itemqty
			added.result <- bm.stockpile.stock[added.itemname]

		case taken = <-bm.stockpile.pickup:
			amount := min(bm.stockpile.stock[taken.itemname], taken.itemqty)
			bm.stockpile.stock[taken.itemname] -= amount
			taken.result <- amount

		case inquiry = <-bm.stockpile.query:
			inquiry.result <- bm.stockpile.stock[inquiry.itemname]
		}
	}
}

// report the makeup of the town's population
func (bm *Burgermeister) listWorkers() {
	for k, v := range bm.workers {
		fmt.Printf("worker %d is %s\n", k, v.name)
	}
	fmt.Printf("\n")
}

// report the contents / status of the town's stockpile
func (bm *Burgermeister) showStockpile() {
	for _, v := range orderedStockpile {
		if _, ok := bm.stockpile.stock[v]; ok {
			fmt.Printf("stock of %s is %d\n", v, bm.stockpile.stock[v])
		}
	}

	fmt.Printf("\n")
}
