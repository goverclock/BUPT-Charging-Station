package scheduler

import (
	"buptcs/data"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

type Scheduler struct {
	mu          sync.Mutex
	stations    []data.Station
	waitingcars []data.Car
}

var sched Scheduler

// initially we have 5 chargin station, 0 cars
func init() {
	sched.mu.Lock()
	defer sched.mu.Unlock()

	go ticker()
	go show_info()

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

// join the car into the waiting queue, so that we can schedule it
func JoinCar(car data.Car) bool {
	sched.mu.Lock()
	defer sched.mu.Unlock()

	// check if the waiting queue is full
	free := data.MAX_WAITING_SLOT - len(sched.waitingcars)
	if free == 0 { // fail
		return false
	} else if free > 0 { // succeed
		car.Stage = data.Waiting
		car.QId = generateQId(car.ChargeMode)
		sched.waitingcars = append(sched.waitingcars, car)
		return true
	}

	log.Fatal("impossible waiting slots:", free)
	return false
}

func CarByUser(u *data.User) (data.Car, error) {
	sched.mu.Lock()
	defer sched.mu.Unlock()
	// 1. find in waitingcars
	for _, c := range sched.waitingcars {
		if c.OwnedBy == u.Uuid {
			return c, nil
		}
	}
	// 2. find in every stations
	for _, st := range sched.stations {
		for _, c := range st.Queue {
			if c.OwnedBy == u.Uuid {
				return c, nil
			}
		}
	}
	return data.Car{}, errors.New("car not found, because user hasn't submit charge")
}

func ticker() {
	for {
		time.Sleep(1 * time.Second)
		sched.mu.Lock()

		// check if any car in waitingcars can move to charing station
		for ci, c := range sched.waitingcars {
			// check if any station is available for the car
			for sti, st := range sched.stations {
				if !st.Available() || st.Mode != c.ChargeMode {
					continue
				}
				// available station
				// generate QId
				// car moves to station
				sched.waitingcars =
					append(sched.waitingcars[:ci], sched.waitingcars[ci+1:]...) // remove from waiting cars
				sched.stations[sti].Join(&c) // join station queue
				break
			}
		}

		sched.mu.Unlock()
	}
}

// assume sched.mu is locked
func generateQId(mode int) string {
	if sched.mu.TryLock() {
		log.Fatal("should have locked sched.mu in generate QId")
	}

	isqidfree := func(s string) bool {
		for _, cs := range sched.waitingcars {
			if cs.QId == s {
				return false
			}
		}
		for _, st := range sched.stations {
			for _, c := range st.Queue {
				if c.QId == s {
					return false
				}
			}
		}
		return true
	}
	for i := 1; ; i++ {
		qid := "T" + strconv.Itoa(i)
		if mode == 1 {
			qid = "F" + strconv.Itoa(i)
		}
		if isqidfree(qid) {
			return qid
		}
	}
}

// for debug only
// Waiting: 0
// Sta1(T):	T1*
// Sta2(T): T2
// Sta3(F): F1*
// Sta4(F):
func show_info() {
	v := os.Getenv("V")
	if v != "" {
		return
	}

	for {
		time.Sleep(1 * time.Second)

		sched.mu.Lock()
		fmt.Printf("Waiting:\t%d\t", len(sched.waitingcars))
		for _, c := range sched.waitingcars {
			fmt.Printf("%v\t", c.QId)
		}
		fmt.Println()
		for _, st := range sched.stations {
			m := 'T'
			if st.Mode == 1 {
				m = 'F'
			}
			fmt.Printf("Sta%d(%c):%d\t", st.Id, m, len(st.Queue))
			for _, c := range st.Queue {
				if c.Stage == data.Charging {
					fmt.Printf("%v*\t", c.QId)
				} else {
					fmt.Printf("%v\t", c.QId)
				}
			}

			fmt.Println()
		}

		fmt.Println()
		sched.mu.Unlock()
	}
}
