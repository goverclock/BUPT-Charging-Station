package data

// Car is created only when a charge request is issued
type Car struct {
	Id      int
	OwnedBy string // user's uuid
	Stage   string // Waiting, Queueing, Charging
	QId     string // F1, F2, T1, T2...
}

func CarById(id int) (c Car, err error) {
	c = Car{}
	err = Db.QueryRow("SELECT id, ownedby, stage, qid FROM cars WHERE id = $1", id).Scan(&c.Id, &c.OwnedBy, &c.Stage, &c.QId)
	return
}
