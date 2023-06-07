package scheduler

import (
	"buptcs/data"
	"buptcs/vtime"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func Assert(b bool, s string) {
	if !b {
		log.Fatal(s)
	}
}

func UserByContext(ctx *gin.Context, user_id int) (user data.User) {
	user_name, ok := ctx.Get("user_name")
	var err error
	if ok {	// using JWT
		user, err = data.UserByName(user_name.(string))
		if err != nil {
			log.Fatal("UserByName: ", user_name)
		}
	} else {	// using user_id
		user, err = data.UserById(user_id)
		if err != nil {
			log.Fatal("UserById: ", user_id)
		}
	}
	return 
}

// assume sched.mu is locked
func generateQId(mode int) string {
	Assert(!sched.mu.TryLock(), "should have locked sched.mu in generate QId")

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
	Assert(!sched.mu.TryLock(), "should have locked sched.mu in getMaxQid")
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

		fmt.Printf("Time:\t%v\n", vtime.Now())
		fmt.Printf("FastInd:\t%v\n", "F"+strconv.Itoa(sched.fast_qind+1))
		fmt.Printf("SlowInd:\t%v\n", "T"+strconv.Itoa(sched.slow_qind+1))
		fmt.Printf("Waiting:\t")
		for _, c := range sched.waitingcars {
			fmt.Printf("%v\t", c.QId)
		}
		fmt.Println()
		fmt.Printf("Temp:\t\t")
		for _, c := range sched.temp_area {
			fmt.Printf("%v\t", c.QId)
		}
		fmt.Println()
		for _, st := range sched.stations {
			m := "T"
			if st.Mode == 1 {
				m = "F"
			}
			if st.GetRunning() {
				m += "1"
			} else {
				m += "0"
			}
			if st.GetIsDown() {
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

		if vtime.ShouldFreeze() {
			fmt.Println("**FREEZED**")
			fmt.Printf("Time:\t%v\n", vtime.Now())
		}
		for vtime.ShouldFreeze() {
			time.Sleep(time.Second)
		}
	}
}

// yuan per kWh
func GetFee() (elec float64, service float64) {
	service = 0.8
	now := vtime.Now()
	// minute := time.Minute
	hour := now.Hour()

	if (hour >= 10 && hour <= 15) || (hour >= 18 && hour <= 21) {
		elec = 1.0
	} else if (hour >= 7 && hour <= 10) || (hour >= 15 && hour <= 18) || (hour >= 21 && hour <= 23) { // 7:00 - 10:00
		elec = 0.7
	} else {
		elec = 0.4
	}
	return
}
