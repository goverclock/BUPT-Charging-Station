package data

import (
	"log"
	"sync"
	"time"
)

const (
	StaOn   int = 0
	StaOff  int = 1
	StaDown int = 2
	StaUp   int = 3
)

type StationReport struct {
	Charge_id         int     `json:"charge_id"`
	Charge_mode       int     `json:"charge_mode"`
	Charge_state      int     `json:"charge_state"`
	Tot_charge_amount float64 `json:"tot_charge_amount"`
	Tot_charge_time   int     `json:"tot_charge_time"`
	Tot_frequency     int     `json:"tot_frequency"`
}

// 2 Fast, 3 Slow
type Station struct {
	Id          int
	Mode        int     // 1 - Fast, 0 - Slow
	Speed       float64 // `Speed` kWh per MINUTE(HOUR in api)
	Queue       []*Car  // the 1st Car can start charge
	ChargeChan  chan float64
	ControlChan chan int

	Running bool
	IsDown  bool

	mu sync.Mutex
}

func NewStation(id int, mode int, speed float64) *Station {
	st := Station{}
	st.Id = id
	st.Mode = mode
	st.Speed = speed
	st.ChargeChan = make(chan float64) // should not buffer too much
	st.ControlChan = make(chan int)
	st.Running = true
	st.IsDown = false
	go st.generateElectricity()
	return &st
}

func (st *Station) GenerateStationReport(start int64, end int64) StationReport {
	st.mu.Lock()
	defer st.mu.Unlock()
	strp := StationReport{}
	strp.Charge_id = st.Id
	strp.Charge_mode = st.Mode
	if st.IsDown {
		strp.Charge_state = 3
	} else if !st.Running {
		strp.Charge_state = 2
	} else if len(st.Queue) != 0 && st.Queue[0].Stage == Charging {
		strp.Charge_state = 1
	} else {
		strp.Charge_state = 0
	}

	// gather information from db
	rp := Report{}
	rows, err := Db.Query("SELECT * FROM reports WHERE charge_id = $1", st.Id)
	if err != nil {
		log.Fatal(err, "Db.Query()")
	}
	defer rows.Close()

	secs := 0
	for rows.Next() {
		err = rows.Scan(&rp.Id, &rp.Num, &rp.Charge_id, &rp.Charge_mode, &rp.Username, &rp.User_id, &rp.Request_charge_amount, &rp.Real_charge_amount, &rp.Charge_time, &rp.Charge_fee, &rp.Service_fee, &rp.Tot_fee, &rp.Step, &rp.Queue_number, &rp.Subtime, &rp.Inlinetime, &rp.Calltime, &rp.Charge_start_time, &rp.Charge_end_time, &rp.Terminate_flag, &rp.Terminate_time, &rp.Failed_flag, &rp.Failed_msg)
		if err != nil {
			log.Fatal(err, "rows.Scan()")
		}
		if rp.Subtime < start || rp.Subtime > end {
			continue
		}
		strp.Tot_charge_amount += rp.Real_charge_amount // total charge amount
		strp.Tot_frequency++                            // total frequency
		if rp.Charge_end_time > rp.Charge_start_time {
			secs += int(rp.Charge_end_time - rp.Charge_start_time)
		}
	}
	strp.Tot_charge_time = secs / 60	// total charge time

	return strp
}

func (st *Station) GetRunning() bool {
	st.mu.Lock()
	defer st.mu.Unlock()
	return st.Running
}

func (st *Station) SetRunning(r bool) {
	st.mu.Lock()
	defer st.mu.Unlock()
	st.Running = r
}

func (st *Station) GetIsDown() bool {
	st.mu.Lock()
	defer st.mu.Unlock()
	return st.IsDown
}

func (st *Station) SetIsDown(d bool) {
	st.mu.Lock()
	defer st.mu.Unlock()
	st.IsDown = d
}

func (st *Station) GetQueue() []*Car {
	st.mu.Lock()
	defer st.mu.Unlock()
	return st.Queue
}

func (st *Station) On() {
	st.ControlChan <- StaOn
}

func (st *Station) Off() {
	st.ControlChan <- StaOff
}

// set as failure
func (st *Station) Down() {
	st.ControlChan <- StaDown
}

// set as non-failure
func (st *Station) Up() {
	st.ControlChan <- StaUp
}

// assumes sched.mu is locked
// put electricity into ChargeChan periodically
func (st *Station) generateElectricity() {
	run := true // only for this function
	up := true

	for {
		time.Sleep(time.Second)

		// first check if should turn on/off the station
		select {
		case ctl := <-st.ControlChan:
			if ctl == StaOn {
				run = true
			} else if ctl == StaOff {
				run = false
			} else if ctl == StaUp {
				up = true
			} else if ctl == StaDown {
				up = false
			}
			st.SetRunning(run)
			st.SetIsDown(!up)
		default:
		}

		// keep trying to send out electricity and
		// simply blocks if no car is receiving electricity
		if up && run {
			select {
			case st.ChargeChan <- st.Speed / 60: // 60 = seconds per minute
			default:
			}
		}
	}
}

func (st *Station) Available() bool {
	return (len(st.Queue) < MAX_STATION_QUEUE) && !st.GetIsDown() && st.GetRunning()
}

func (st *Station) Join(c *Car) {
	if !st.Available() {
		log.Fatal("station is not available!")
	}
	st.Queue = append(st.Queue, c)
}

func (st *Station) Leave(qid string) {
	if len(st.Queue) != 0 && st.Queue[0].QId == qid {
		st.Queue = st.Queue[1:]
	} else {
		log.Fatal("st.Leave()")
	}
}

// when the station fails
func (st *Station) LeaveAll() []*Car {
	ret := st.Queue
	st.Queue = nil
	return ret
}

// time needed for car before finishing charging
func (st *Station) WaitingTimeForCar(c Car) float64 {
	ret := c.ChargeAmount / st.Speed
	for _, c := range st.Queue {
		ret += c.ChargeAmount / st.Speed
	}
	return ret
}
