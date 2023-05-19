package scheduler

import (
	"buptcs/data"
	"log"
	"strconv"
	"sync"
	"time"
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

	go ticker()
	sched.stations = append(sched.stations, data.Station{
		Id:   0,
		Mode: 1,
	})
	sched.stations = append(sched.stations, data.Station{
		Id:   1,
		Mode: 1,
	})
	sched.stations = append(sched.stations, data.Station{
		Id:   2,
		Mode: 0,
	})
	sched.stations = append(sched.stations, data.Station{
		Id:   3,
		Mode: 0,
	})
	sched.stations = append(sched.stations, data.Station{
		Id:   4,
		Mode: 0,
	})
}

// join the car into the queue, so that we can schedule it
func JoinCar(car data.Car) bool {
	sched.mu.Lock()
	defer sched.mu.Unlock()

	// check if the waiting queue is full
	waiting_slots := data.MAX_WAITING_SLOT
	for _, c := range sched.cars {
		if c.Stage == "Waiting" {
			waiting_slots--
		}
	}
	if waiting_slots == 0 { // fail
		return false
	} else if waiting_slots > 0 { // succeed
		car.Stage = "Waiting"
		sched.cars = append(sched.cars, car)
		return true
	}
	log.Fatal("impossible waiting_slots:", waiting_slots)
	return false
}

func ticker() {
	for {
		time.Sleep(1 * time.Second)
		sched.mu.Lock()
		// check if any car can move from waiting area to charing station's slots
		for ci, c := range sched.cars {
			if c.Stage != "Waiting" {
				continue
			}
			// check if any station is available
			for _, st := range sched.stations {
				if !st.Available() || st.Mode != c.ChargeMode {
					continue
				}
				// available station
				// generate QId
				qidfree := func(s string) bool {
					for _, cs := range sched.cars {
						if cs.QId == s {
							return false
						}
					}
					return true
				}
				for i := 1; ;i++ {
					qid := "T" + strconv.Itoa(i)
					if sched.cars[ci].ChargeMode == 1 {
						qid = "F" + strconv.Itoa(i)
					}
					if qidfree(qid) {
						sched.cars[ci].QId = qid
						break
					}
				}
				// car moves to station
				sched.cars[ci].Stage = "Queueing"
				st.Join(&c)
				break
			}
		}

		sched.mu.Unlock()
	}
}
