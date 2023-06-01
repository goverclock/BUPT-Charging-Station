package scheduler

import (
	"buptcs/data"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

// assume sched.mu is locked
func generateQId(mode int) string {
	if sched.mu.TryLock(){
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
	// return data.MAX_WAITING_SLOT + data.MAX_STATION_QUEUE*len(sched.stations)
	return 100
}

// for debug only
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
		fmt.Printf("Temp:\t")
		for _, c := range sched.temp_area {
			fmt.Printf("%v\t", c.QId)
		}
		fmt.Println()
		for _, st := range sched.stations {
			m := "T"
			if st.Mode == 1 {
				m = "F"
			}
			if st.Running {
				m += "1"
			} else {
				m += "0"
			}
			if st.IsDown {
				m += "0"
			} else {
				m += "1"
			}
			fmt.Printf("Sta%d(%s):%d\t", st.Id, m, len(st.Queue))
			for _, c := range st.Queue {
				if c.Stage == data.Charging {
					fmt.Printf("%v*\t", c.QId)
				} else {
					fmt.Printf("%v\t", c.QId)
				}
			}
			fmt.Println()
		}
		fmt.Println(sched.ongoing_reports)

		fmt.Println()
		sched.mu.Unlock()
	}
}

// yuan per kWh
func getFee() (elec float64, service float64)  {
	service = 0.8
	now := time.Now()
	// minute := time.Minute
	hour := now.Hour()
	
	if (hour >= 10 && hour <= 15) || (hour >= 18 && hour <= 21) {
		elec = 1.0
	} else if (hour >= 7 && hour <= 10) || (hour >= 15 && hour <= 18) || (hour >= 21 && hour <= 23) {	// 7:00 - 10:00
		elec = 0.7
	} else {
		elec = 0.4
	}
	return
}
