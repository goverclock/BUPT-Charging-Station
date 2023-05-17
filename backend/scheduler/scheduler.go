package scheduler

import (
	"buptcs/data"
	"sync"
)

type Scheduler struct {
	mu sync.Mutex
}

var Sched Scheduler

func init() {

}

// join the car into the queue, so that we can schedule it
func (sch *Scheduler)JoinCar(car data.Car) {

}

