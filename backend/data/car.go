package data

// Car is created only when a charge request is issued
type Car struct {
	Id      int	// not used
	OwnedBy string // user's uuid
	Stage   string // Waiting, Queueing, Charging
	QId     string // F1, F2, T1, T2...
}
