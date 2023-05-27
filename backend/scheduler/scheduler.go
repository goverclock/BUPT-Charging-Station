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
	fast_qind   int // QId of the next fast car to schdule to charging area
	slow_qind   int
}

var sched Scheduler

// initially we have 5 chargin station, 0 cars
func init() {
	sched.mu.Lock()
	defer sched.mu.Unlock()

	go ticker()
	go show_info()

	sched.stations = append(sched.stations, data.Station{
		Id:    0,
		Mode:  1,
		Speed: 30,
	})
	sched.stations = append(sched.stations, data.Station{
		Id:    1,
		Mode:  1,
		Speed: 30,
	})
	sched.stations = append(sched.stations, data.Station{
		Id:    2,
		Mode:  0,
		Speed: 7,
	})
	sched.stations = append(sched.stations, data.Station{
		Id:    3,
		Mode:  0,
		Speed: 7,
	})
	sched.stations = append(sched.stations, data.Station{
		Id:    4,
		Mode:  0,
		Speed: 7,
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

		// try to schdule the next fast car
		schduleFast()
		// try to schedule the next slow car
		scheduleSlow()

		sched.mu.Unlock()
	}
}

func schduleFast() {
	for ci, c := range sched.waitingcars {
		if c.QId == "F"+strconv.Itoa(sched.fast_qind+1) {
			// the car exists, look for a station with min wait time for the car
			min_wait := -1.0
			min_wait_sti := -1
			for sti, st := range sched.stations {
				if !st.Available() || st.Mode != 1 {	// fast station
					continue
				}
				if min_wait < 0 || st.WaitingTimeForCar(c) < min_wait {
					min_wait = st.WaitingTimeForCar(c)
					min_wait_sti = sti
				}
			}

			// available station
			// car moves to station
			if min_wait_sti != -1 {
				sched.waitingcars =
					append(sched.waitingcars[:ci], sched.waitingcars[ci+1:]...) // remove from waiting cars
				sched.stations[min_wait_sti].Join(&c) // join station queue
				sched.fast_qind++
			}

			break
		}
	}
}

func scheduleSlow() {
	for ci, c := range sched.waitingcars {
		if c.QId == "T"+strconv.Itoa(sched.slow_qind+1) {
			// the car exists, look for a station with min wait time for the car
			min_wait := -1.0
			min_wait_sti := -1
			for sti, st := range sched.stations {
				if !st.Available() || st.Mode != 0 {	// slow station
					continue
				}
				if min_wait < 0 || st.WaitingTimeForCar(c) < min_wait {
					min_wait = st.WaitingTimeForCar(c)
					min_wait_sti = sti
				}
			}

			// available station
			// car moves to station
			if min_wait_sti != -1 {
				sched.waitingcars =
					append(sched.waitingcars[:ci], sched.waitingcars[ci+1:]...) // remove from waiting cars
				sched.stations[min_wait_sti].Join(&c) // join station queue
				sched.slow_qind++
			}

			break
		}
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
	if mode == 1 { // fast
		for i := sched.fast_qind; ; i = (i + 1) % getMaxQId() {
			qid := "F" + strconv.Itoa(i+1)
			if isqidfree(qid) {
				return qid
			}
		}
	} else { // slow
		for i := sched.slow_qind; ; i = (i + 1) % getMaxQId() {
			qid := "T" + strconv.Itoa(i+1) // slow
			if isqidfree(qid) {
				return qid
			}
		}
	}
}

// assume sched.mu is locked
func getMaxQId() int {
	if sched.mu.TryLock() {
		log.Fatal("should have locked sched.mu in getMaxQid")
	}
	return data.MAX_WAITING_SLOT + data.MAX_STATION_QUEUE*len(sched.stations)
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

		fmt.Printf("FastInd:\t%v\n", "F"+strconv.Itoa(sched.fast_qind+1))
		fmt.Printf("SlowInd:\t%v\n", "T"+strconv.Itoa(sched.slow_qind+1))
		fmt.Printf("Waiting:\t")
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
