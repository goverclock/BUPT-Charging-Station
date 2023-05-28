package data

import "log"

// 2 Fast, 3 Slow
type Station struct {
	Id    int
	Mode  int // 1 - Fast, 0 - Slow
	Speed float64
	Queue []*Car // the 1st Car can start charge

	Running bool
	Failure bool
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
