package data

type Station struct {
	Id     int
	Mode   string // Fast, Slow
	UsedBy string // F1, F2, T1, T2...
	Slot1  string // F1, F2, T1, T2...
	Slot2  string
}

func StationById(id int) (st Station, err error) {
	st = Station{}
	err = Db.QueryRow("SELECT id, mode, usedby, slot1, slot2 FROM stations WHERE id = $1", id).Scan(&st.Id, &st.Mode, &st.UsedBy, &st.Slot1, &st.Slot2)
	return
}
