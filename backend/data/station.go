package data

// 2 Fast, 3 Slow
type Station struct {
	Id     int
	Mode   string // Fast, Slow
	UsedBy string // F1, F2, T1, T2...
	Slot1  string // F1, F2, T1, T2...
	Slot2  string
}


