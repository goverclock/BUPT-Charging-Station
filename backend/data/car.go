package data

type CarStage int
const (
	Waiting CarStage = iota
	Queueing
	Charging
)

// Car is created only when a charge request is issued
type Car struct {
	OwnedBy string // user's uuid
	Stage   CarStage // Waiting, Queueing, Charging
	QId     string // F1, F2, T1, T2...
	ChargeMode int // 1 - Fast, 0 - Slow
	ChargeAmount float64
}