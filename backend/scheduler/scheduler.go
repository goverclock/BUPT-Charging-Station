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
		car.Stage = "Waiting"
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
			// check if any station is available
			for sti, st := range sched.stations {
				if !st.Available() || st.Mode != c.ChargeMode {
					continue
				}
				// available station
				// generate QId
				qidfree := func(s string) bool {
					for _, cs := range sched.waitingcars {
						if cs.QId == s {
							return false
						}
					}
					return true
				}
				for i := 1; ; i++ {
					qid := "T" + strconv.Itoa(i)
					if sched.waitingcars[ci].ChargeMode == 1 {
						qid = "F" + strconv.Itoa(i)
					}
					if qidfree(qid) {
						sched.waitingcars[ci].QId = qid
						break
					}
				}
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

// for debug only
// Waiting: 0
// Sta1(T):	T1*
// Sta2(T): T2
// Sta3(F): F1*
// Sta4(F): 
// 
func show_info() {
	v := os.Getenv("V")
	if v != "" {
		return
	}
	
	for {
		time.Sleep(1 * time.Second)
		
		sched.mu.Lock()
		fmt.Printf("Waiting:\t%d\n", len(sched.waitingcars))
		for _, st := range sched.stations {
			m := 'T'
			if st.Mode == 1 {
				m = 'F'
			}
			fmt.Printf("Sta%d(%c):\t%d", st.Id, m, len(st.Queue))
			// for _, c := range st.Queue {
			// 	fmt.Printf("%s\t" ,c.QId)
			// }

			fmt.Println()
		}
		
		fmt.Println()
		sched.mu.Unlock()
	}
}
