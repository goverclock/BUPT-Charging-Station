package data

// Car is created only when a charge request is issued
type Car struct {
	Id      int
	OwnedBy string // user's uuid
	Stage   string // Waiting, Queueing, Charging
	QId     string // F1, F2, T1, T2...
}

var Cars []Car
