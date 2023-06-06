package scheduler

import (
	"buptcs/data"
	"buptcs/vtime"
	"errors"
	"log"
	"strconv"
	"sync"
	"time"
)

type Scheduler struct {
	mu          sync.Mutex
	stations    []*data.Station
	waitingcars []*data.Car
	temp_area   []*data.Car
	fast_qind   int // QId of the next fast car to schdule to charging area
	slow_qind   int

	ongoing_reports []*data.Report // every user should have at most 1 ongoing report
	fault_schedule  int            // 优先级调度:0/时间顺序调度:1
}

var sched Scheduler

// initially we have 5 chargin station, 0 cars
func init() {
	sched.mu.Lock()
	defer sched.mu.Unlock()

	go ticker()
	go show_info()

	// create stations
	stid := 0
	for ; stid < data.FAST_STATION_COUNT; stid++ {
		st := data.NewStation(stid, 1, 30)
		sched.stations = append(sched.stations, st)
	}
	for ; stid < data.SLOW_STATION_COUNT+data.FAST_STATION_COUNT; stid++ {
		st := data.NewStation(stid, 0, 7)
		sched.stations = append(sched.stations, st)
	}
}

// join the car into the waiting queue, so that we can schedule it
func JoinCar(user data.User, car *data.Car) bool {
	sched.mu.Lock()
	defer sched.mu.Unlock()

	// create an ongoing report
	rp := newOngoingReport(user)

	// check if the waiting queue is full
	free := data.MAX_WAITING_SLOT - len(sched.waitingcars)
	if free == 0 { // fail
		rp.Failed_flag = true
		rp.Failed_msg = "等待队列已满,无法加入"
		archiveOngoingReport(rp) // charge end, save report(which failed)
		return false
	} else if free > 0 { // succeed
		car.Stage = data.Waiting
		car.QId = generateQId(car.ChargeMode)
		sched.waitingcars = append(sched.waitingcars, car)
		rp.Charge_mode = car.ChargeMode // update report
		rp.Request_charge_amount = car.ChargeAmount
		rp.Step = data.StepInline
		rp.Queue_number = car.QId
		return true
	}

	log.Println("impossible waiting slots:", free)
	return false
}

func CarByUser(u data.User) (*data.Car, error) {
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
	return &data.Car{}, errors.New("car not found, because user hasn't submit charge")
}

// if no report found, returns a dumb value with Num = 0
func OngoingCopyByUser(u data.User) data.Report {
	sched.mu.Lock()
	defer sched.mu.Unlock()
	for _, rp := range sched.ongoing_reports {
		if rp.Username == u.Name {
			return *rp
		}
	}
	return data.Report{Num: 0}
}

// assume sched.mu is locked
func stationById(id int) *data.Station {
	Assert(!sched.mu.TryLock(), "should have locked sched.mu in stationById")
	for _, st := range sched.stations {
		if st.Id == id {
			return st
		}
	}
	log.Println("can't find station with id ", id)
	return nil
}

func StartChargeCar(c *data.Car) error {
	sched.mu.Lock()
	defer sched.mu.Unlock()
	if c.Stage == data.Waiting {
		return errors.New("car isn't queueing")
	} else if c.Stage == data.Charging {
		return errors.New("car is alreay charging")
	}
	// check if the car is in a station's 1st slot
	for _, st := range sched.stations {
		if len(st.Queue) > 0 && st.Queue[0].QId == c.QId {
			// start charge
			st.Queue[0].Stage = data.Charging

			// update report
			rp := ongoingReportByUser(data.UserByUUId(c.OwnedBy))
			rp.Step = data.StepCharge
			rp.Charge_start_time = vtime.Now().Unix()
			return nil
		}
	}
	return errors.New("no such car")
}

// 0 if car is not in waiting area
func WaitCountByCar(c *data.Car) (int, error) {
	sched.mu.Lock()
	defer sched.mu.Unlock()

	if c.Stage != data.Waiting {
		return 0, nil
	}
	qid := c.QId
	wc := 0
	if c.ChargeMode == 1 {
		for {
			if "F"+strconv.Itoa(((sched.fast_qind+wc)%getMaxQId())+1) == qid {
				return wc, nil
			}
			wc++
			if wc > getMaxQId() {
				return 0, errors.New("waiting_count out of range")
			}
		}
	} else {
		for {
			if "T"+strconv.Itoa(((sched.slow_qind+wc)%getMaxQId())+1) == qid {
				return wc, nil
			}
			wc++
			if wc > getMaxQId() {
				return 0, errors.New("waiting_count out of range")
			}
		}
	}
}

func StationById(stid int) (*data.Station, error) {
	sched.mu.Lock()
	defer sched.mu.Unlock()
	for _, st := range sched.stations {
		if st.Id != stid {
			continue
		}
		return st, nil
	}
	return &data.Station{}, errors.New("no such station")
}

func GetFautlSchedule() int {
	sched.mu.Lock()
	defer sched.mu.Unlock()
	return sched.fault_schedule
}

func ChangeSettings(was int, cql int, cs int, fs int) {
	sched.mu.Lock()
	defer sched.mu.Unlock()
	data.MAX_WAITING_SLOT = was
	data.MAX_STATION_QUEUE = cql
	data.CALL_SCHEDULE = cs
	sched.fault_schedule = fs
}

// on/off station
func SwitchStation(stid int, is_on bool) {
	sched.mu.Lock()
	defer sched.mu.Unlock()
	for _, st := range sched.stations {
		if st.Id == stid {
			if is_on {
				st.On()
			} else {
				st.Off()
			}
			break
		}
	}
}

func AllStationReports(start int64, end int64) (strps []data.StationReport) {
	sched.mu.Lock()
	defer sched.mu.Unlock()
	for _, st := range sched.stations {
		strps = append(strps, st.GenerateStationReport(start, end))
	}
	return
}

func SwitchBrokenStation(stid int, is_fail bool) {
	sched.mu.Lock()
	defer sched.mu.Unlock()
	for _, st := range sched.stations {
		if st.Id == stid {
			if is_fail {
				st.Down()
			} else {
				st.Up()
			}
			break
		}
	}
}

// user must have submitted the charge
// and should not be charging
func CancelCharge(u data.User) bool {
	sched.mu.Lock()
	defer sched.mu.Unlock()

	rp := ongoingReportByUser(u)
	if rp == nil || rp.Step == data.StepCharge {
		return false
	}
	// cancel the charge(report)
	now := vtime.Now().Unix()
	rp.Terminate_flag = true
	rp.Terminate_time = now
	rp.Step = data.StepFinish
	archiveOngoingReport(rp)

	// remove user'car from waiting area/station's queue
	removeCar(u)

	return true
}

// user must be charging
func EndCharge(u data.User) bool {
	sched.mu.Lock()
	defer sched.mu.Unlock()

	rp := ongoingReportByUser(u)
	if rp == nil || rp.Step != data.StepCharge {
		return false
	}

	// end the charge(report)
	now := vtime.Now().Unix()
	rp.Step = data.StepFinish
	rp.Charge_end_time = now
	rp.Terminate_flag = true
	archiveOngoingReport(rp)

	removeCar(u)

	return true
}

// user must be in waiting area
// if mode(fast/slow) changes, generate new QId
func ChangeCharge(u data.User, mode int, amount float64) bool {
	sched.mu.Lock()
	defer sched.mu.Unlock()

	rp := ongoingReportByUser(u)
	if rp == nil || rp.Step != data.StepInline {
		return false
	}

	// changing ChargeAmount is straightforward
	rp.Request_charge_amount = amount
	// if mode changes, there's some extra work
	if rp.Charge_mode != mode {
		removeCar(u) // first remove from waiting area
		rp.Charge_mode = mode
		rp.Queue_number = generateQId(mode)
		changeCar(u, mode, amount, rp.Queue_number)
	}

	return true
}

func ticker() {
	for {
		time.Sleep(1 * time.Second)
		sched.mu.Lock()

		if checkFault() { // now should apply fault schedule
			scheduleFaultFast() // both in QId order
			scheduleFaultSlow()
		} else { // no fault occurs
			// try to schdule the next fast car
			schduleFast()
			// try to schedule the next slow car
			scheduleSlow()
		}

		// cars in stations 1st slot should turn Stage from
		// Queueing to Called
		scheduleCall()

		// actually charge happens here
		updateOngoingReports()

		sched.mu.Unlock()
	}
}

// returns true if some car is in temp_area
// or some car is in a failed station
// should enter fault schedule if returns true
// also moves cars to temp area if necessary
func checkFault() bool {
	// check if some car is in a failed station
	for _, st := range sched.stations {
		if !st.GetIsDown() || len(st.GetQueue()) == 0 {
			continue
		}
		cars := st.LeaveAll()          // move all cars in failed station to temp area
		if sched.fault_schedule == 1 { // involves other stations
			for _, ost := range sched.stations {
				if ost.Mode != st.Mode {
					continue
				}
				for _, oc := range ost.GetQueue() {
					if oc.Stage == data.Queueing {
						cars = append(cars, ost.Leave(oc.QId))
					}
				}
			}
		}
		for _, c := range cars { // and end these reports
			user := data.UserByUUId(c.OwnedBy)
			rp := ongoingReportByUser(user)
			rp.Failed_flag = true
			rp.Failed_msg = "充电桩故障"
			rp.Step = data.StepFinish
			already_charged_amount := rp.Real_charge_amount
			if rp.Charge_start_time != 0 {
				rp.Charge_end_time = vtime.Now().Unix()
			}
			archiveOngoingReport(rp)
			// start a new report for cars in temp area
			nrp := newOngoingReport(user)
			nrp.Charge_mode = c.ChargeMode
			nrp.Request_charge_amount = c.ChargeAmount - already_charged_amount
			nrp.Inlinetime = nrp.Subtime
			nrp.Step = data.StepInline
			nrp.Queue_number = c.QId
		}
		sched.temp_area = append(sched.temp_area, cars...)

		break // deal with 1 failure each call is enough
	}
	return len(sched.temp_area) != 0
}

func schduleFast() {
	for ci, c := range sched.waitingcars {
		if c.QId == "F"+strconv.Itoa(sched.fast_qind+1) {
			// the car exists, look for a station with min wait time for the car
			min_wait := -1.0
			min_wait_sti := -1
			for sti, st := range sched.stations {
				if !st.Available() || st.Mode != 1 { // fast station
					continue
				}
				if min_wait < 0 || st.WaitingTimeForCar(*c) < min_wait {
					min_wait = st.WaitingTimeForCar(*c)
					min_wait_sti = sti
				}
			}

			// available station
			// car moves to station sti
			if min_wait_sti != -1 {
				sched.waitingcars =
					append(sched.waitingcars[:ci], sched.waitingcars[ci+1:]...) // remove from waiting cars
					// join station queue
				c.Stage = data.Queueing // i.e. Inline

				// update report
				rp := ongoingReportByUser(data.UserByUUId(c.OwnedBy))
				rp.Inlinetime = vtime.Now().Unix()
				rp.Charge_id = min_wait_sti

				sched.stations[min_wait_sti].Join(c)
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
				if !st.Available() || st.Mode != 0 { // slow station
					continue
				}
				if min_wait < 0 || st.WaitingTimeForCar(*c) < min_wait {
					min_wait = st.WaitingTimeForCar(*c)
					min_wait_sti = sti
				}
			}

			// available station
			// car moves to station
			if min_wait_sti != -1 {
				sched.waitingcars =
					append(sched.waitingcars[:ci], sched.waitingcars[ci+1:]...) // remove from waiting cars
					// join station queue
				c.Stage = data.Queueing

				// update report
				rp := ongoingReportByUser(data.UserByUUId(c.OwnedBy))
				rp.Inlinetime = vtime.Now().Unix()
				rp.Charge_id = min_wait_sti

				sched.stations[min_wait_sti].Join(c)
				sched.slow_qind++
			}

			break
		}
	}
}

func scheduleFaultFast() {
	is_qid_in_temp := func(qid string) bool {
		for _, c := range sched.temp_area {
			if c.QId == qid {
				return true
			}
		}
		return false
	}
	qid := "F1"
	for i := 2; !is_qid_in_temp(qid) && i < getMaxQId(); i++ {
		qid = "F" + strconv.Itoa(i)
	}

	for ci, c := range sched.temp_area {
		if c.QId != qid {
			continue
		}
		// look for a station with min wait time for the car
		min_wait := -1.0
		min_wait_sti := -1
		for sti, st := range sched.stations {
			if !st.Available() || st.Mode != 1 { // fast station
				continue
			}
			if min_wait < 0 || st.WaitingTimeForCar(*c) < min_wait {
				min_wait = st.WaitingTimeForCar(*c)
				min_wait_sti = sti
			}
		}
		// available station
		// car moves to station
		if min_wait_sti != -1 {
			sched.temp_area =
				append(sched.temp_area[:ci], sched.temp_area[ci+1:]...) // remove from waiting cars
				// join station queue
			c.Stage = data.Queueing
			// update report
			rp := ongoingReportByUser(data.UserByUUId(c.OwnedBy))
			rp.Inlinetime = vtime.Now().Unix()
			rp.Charge_id = min_wait_sti
			sched.stations[min_wait_sti].Join(c)
		}
	}
}

func scheduleFaultSlow() {
	is_qid_in_temp := func(qid string) bool {
		for _, c := range sched.temp_area {
			if c.QId == qid {
				return true
			}
		}
		return false
	}
	qid := "T1"
	for i := 2; !is_qid_in_temp(qid) && i < getMaxQId(); i++ {
		qid = "T" + strconv.Itoa(i)
	}

	for ci, c := range sched.temp_area {
		if c.QId != qid {
			continue
		}
		// look for a station with min wait time for the car
		min_wait := -1.0
		min_wait_sti := -1
		for sti, st := range sched.stations {
			if !st.Available() || st.Mode != 0 { // fast station
				continue
			}
			if min_wait < 0 || st.WaitingTimeForCar(*c) < min_wait {
				min_wait = st.WaitingTimeForCar(*c)
				min_wait_sti = sti
			}
		}
		// available station
		// car moves to station
		if min_wait_sti != -1 {
			sched.temp_area =
				append(sched.temp_area[:ci], sched.temp_area[ci+1:]...) // remove from waiting cars
				// join station queue
			c.Stage = data.Queueing
			// update report
			rp := ongoingReportByUser(data.UserByUUId(c.OwnedBy))
			rp.Inlinetime = vtime.Now().Unix()
			rp.Charge_id = min_wait_sti
			sched.stations[min_wait_sti].Join(c)
		}
	}
}

func scheduleCall() {
	for _, st := range sched.stations {
		if len(st.Queue) == 0 {
			continue
		}
		car := st.Queue[0]
		if car.Stage != data.Queueing {
			continue
		}
		car.Stage = data.Called
		// update report
		user := data.UserByUUId(car.OwnedBy)
		rp := ongoingReportByUser(user)
		rp.Calltime = vtime.Now().Unix()
		rp.Step = data.StepCall
	}

	// if a charging car moved to temp_area due to station failure
	for _, tc := range sched.temp_area {
		if tc.Stage != data.Called {
			continue
		}
		tc.Stage = data.Queueing
		user := data.UserByUUId(tc.OwnedBy)
		rp := ongoingReportByUser(user)
		if rp == nil {
			log.Println("no ongoing report for user!")
			continue
		}
		rp.Calltime = 0
		rp.Step = data.StepInline
	}
}

func removeCar(u data.User) {
	// look for user's car in waiting area
	for ci, c := range sched.waitingcars {
		if c.OwnedBy == u.Uuid {
			sched.waitingcars = append(sched.waitingcars[:ci], sched.waitingcars[:ci+1]...)
			return
		}
	}
	// look for the car in station's queue
	for _, st := range sched.stations {
		for _, c := range st.Queue {
			if c.OwnedBy == u.Uuid {
				st.Leave(c.QId)
				return
			}
		}
	}
}

func changeCar(u data.User, mode int, amount float64, qid string) {
	// look for user's car in waiting area
	for _, c := range sched.waitingcars {
		if c.OwnedBy == u.Uuid {
			c.ChargeMode = mode
			c.ChargeAmount = amount
			c.QId = qid
			return
		}
	}
	// look for the car in station's queue
	for _, st := range sched.stations {
		for _, c := range st.Queue {
			if c.OwnedBy == u.Uuid {
				c.ChargeMode = mode
				c.ChargeAmount = amount
				c.QId = qid
				return
			}
		}
	}
}
