package data

import "log"

// 2 Fast, 3 Slow
type Station struct {
	Id    int
	Mode  int   // 1 - Fast, 0 - Slow
	Queue []Car // the 1st Car can start charge
}

func (st *Station) Available() bool {
	return len(st.Queue) < MAX_STATION_QUEUE
}

func (st *Station) Join(c *Car) {
	st.Queue = append(st.Queue, *c)
	if st.Queue[0] == *c {
		c.Stage = "Charging"
	}
	log.Println(st.Id, "car JOINED", c)
}
