package scheduler

import (
	"buptcs/data"
	"log"
	"sync"
)

type Scheduler struct {
	mu       sync.Mutex
	stations []data.Station
	cars     []data.Car
}

var sched Scheduler

// initially we have 5 chargin station, 0 cars
func init() {
	sched.mu.Lock()
	defer sched.mu.Unlock()
	sched.stations = append(sched.stations, data.Station{
		Id:   0,
		Mode: "Fast",
	})
	sched.stations = append(sched.stations, data.Station{
		Id:   1,
		Mode: "Fast",
	})
	sched.stations = append(sched.stations, data.Station{
		Id:   2,
		Mode: "Slow",
	})
	sched.stations = append(sched.stations, data.Station{
		Id:   3,
		Mode: "Slow",
	})
	sched.stations = append(sched.stations, data.Station{
		Id:   4,
		Mode: "Slow",
	})
}

// join the car into the queue, so that we can schedule it
func JoinCar(car data.Car) bool {
	sched.mu.Lock()
	defer sched.mu.Unlock()

	waiting_slots := data.MAX_WAITING_SLOT
	for _, c := range  sched.cars {
		if c.Stage == "Waiting" {
			waiting_slots--
		}
	}
	if waiting_slots == 0 {
		return false
	} else if waiting_slots > 0 {
		car.Stage = "Waiting"
		sched.cars = append(sched.cars, car)
		return true
	}
	log.Fatal("impossible waiting_slots:", waiting_slots)
	return false
}
