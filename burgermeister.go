package main

import (
	"fmt"
)

const startingBread = 1000
const minimumBreadToRecruitBreadWorkers = startingBread * .95
const minimumBreadToRecruitGeneralWorkers = startingBread * 2

const breadForAMeal = 4

type Burgermeister struct {
	stockpile        *Stockpile
	workers          []*Worker
	recruitmentRules []*RecruitmentRule
}

func (bm *Burgermeister) initializeBurg() {
	bm.stockpile = &Stockpile{}

	bm.stockpile.stock = make(map[string]int, 0)
	bm.stockpile.stock["bread"] = startingBread

	bm.stockpile.dropoff = make(chan Stockupdate)
	bm.stockpile.pickup = make(chan Stockupdate)
	bm.stockpile.query = make(chan string)

	bm.workers = make([]*Worker, 0, 0)

	bm.recruitmentRules = make([]*RecruitmentRule, 0, 0)
	bm.recruitmentRules = append(bm.recruitmentRules, &RecruitmentRule{
		item:      "wheat",
		cmp:       "<",
		threshold: 10,
		worker: &Worker{
			name: "wheat farmer",
			job:  &Wheatfarmer{},
		},
	})
	bm.recruitmentRules = append(bm.recruitmentRules, &RecruitmentRule{
		item:      "flour",
		cmp:       "<",
		threshold: 10,
		worker: &Worker{
			name: "miller",
			job:  &Miller{},
		},
	})
	bm.recruitmentRules = append(bm.recruitmentRules, &RecruitmentRule{
		item:      "bread",
		cmp:       "<",
		threshold: 1000,
		worker: &Worker{
			name: "baker",
			job:  &Baker{},
		},
	})
}

func (bm *Burgermeister) recruitWorkers() {
	if bm.stockpile.stock["bread"] < minimumBreadToRecruitBreadWorkers {
		fmt.Printf(">>>>  not enough bread to recruit new bread workers !\n")
		return
	}

	for _, v := range bm.recruitmentRules {
		if v.cmp == "<" {
			if bm.stockpile.stock[v.item] < v.threshold {
				bm.recruitWorker(v.worker)
			}
		}

		if v.cmp == "==" {
			if bm.stockpile.stock[v.item] == v.threshold {
				bm.recruitWorker(v.worker)
			}
		}

		if v.cmp == ">" {
			if bm.stockpile.stock[v.item] > v.threshold {
				bm.recruitWorker(v.worker)
			}
		}
	}
}

func (bm *Burgermeister) feedWorkers() {
	for _, v := range bm.workers {
		bm.feedWorker(v)
	}
}

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

func (bm *Burgermeister) recruitWorker(worker *Worker) {
	fmt.Printf("Burgermeister.recruitWorker : %s\n", worker.name)

	bm.workers = append(bm.workers, worker)
	go worker.job.goToWork(bm.stockpile)
}

func (bm *Burgermeister) updateStockpile() {
	var added Stockupdate
	var taken Stockupdate

	for {
		select {
		case added = <-bm.stockpile.dropoff:
			bm.stockpile.stock[added.itemname] += added.itemqty
			added.result <- bm.stockpile.stock[added.itemname]

		case taken = <-bm.stockpile.pickup:
			amount := min(bm.stockpile.stock[taken.itemname], taken.itemqty)
			bm.stockpile.stock[taken.itemname] -= amount
			taken.result <- amount
		}
	}
}

func (bm *Burgermeister) listWorkers() {
	for k, v := range bm.workers {
		fmt.Printf("worker %d is %s\n", k, v.name)
	}
	fmt.Printf("\n")
}

func (bm *Burgermeister) showStockpile() {
	for k, v := range bm.stockpile.stock {
		fmt.Printf("stock of %s is %d\n", k, v)
	}
	fmt.Printf("\n")
}

type RecruitmentRule struct {
	item      string
	cmp       string
	threshold int
	worker    *Worker
}
