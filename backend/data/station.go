package data

import (
	"log"
	"time"
)

// 2 Fast, 3 Slow
type Station struct {
	Id         int
	Mode       int     // 1 - Fast, 0 - Slow
	Speed      float64 // `Speed` kWh per MINUTE(HOUR in api)
	Queue      []*Car  // the 1st Car can start charge
	ChargeChan chan float64
	ControlChan chan bool

	Running bool
	Failure bool
}

func NewStation(id int, mode int, speed float64) *Station {
	st := Station{}
	st.Id = id
	st.Mode = mode
	st.Speed = speed
	st.ChargeChan = make(chan float64) // should not buffer too much
	go st.generateElectricity()
	return &st
}

// put electricity into ChargeChan periodically
func (st *Station) generateElectricity() {
	for {
		time.Sleep(time.Second)

		// keep trying to send out electricity and
		// simply blocks if no car is receiving electricity
		st.ChargeChan <- st.Speed / 60 // 60 = seconds per minute
		// log.Println(st.Id, " generated ", st.Speed /60)
	 }
}

func (st *Station) Available() bool {
	return len(st.Queue) < MAX_STATION_QUEUE
}

func (st *Station) Join(c *Car) {
	if !st.Available() {
		log.Fatal("station is not available!")
	}
	st.Queue = append(st.Queue, c)
	// log.Println(st.Id, "car JOINED", c)
}

// time needed for car before finishing charging
func (st *Station) WaitingTimeForCar(c Car) float64 {
	ret := c.ChargeAmount / st.Speed
	for _, c := range st.Queue {
		ret += c.ChargeAmount / st.Speed
	}
	return ret
}
