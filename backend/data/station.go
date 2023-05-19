package data

// 2 Fast, 3 Slow
type Station struct {
	Id     int
	Mode   int    // 1 - Fast, 0 - Slow
	Queue  []string	// QIds
}

func (st *Station) Available() bool {
	return len(st.Queue) < MAX_STATION_QUEUE
}

func (st *Station) Join(c *Car) {
	st.Queue = append(st.Queue, c.QId)	
	if st.Queue[0] == c.QId {
		c.Stage = "Charging"
	}
	// log.Println("car JOINED", c)
}
