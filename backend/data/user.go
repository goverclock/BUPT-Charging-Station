package data

import "time"

type User struct {
	Id              int
	Name            string
	Password        string
	Balance         float64
	BatteryCapacity float64
}

type Session struct {
	Id        int
	Uuid      string
	Name      string
	UserId    int
	CreatedAt time.Time
}

// check if session is valid in the database
// TODO: legacy
func (session *Session) Check() (valid bool, err error) {
	err = Db.QueryRow("SELECT id, uuid, name, user_id, created_at FROM sessions WHERE uuid = $1", session.Uuid).
		Scan(&session.Id, &session.Uuid, &session.Name, &session.UserId, &session.CreatedAt)
	if err != nil {
		valid = false
		return
	}
	if session.Id != 0 {
		valid = true
	}
	return
}

// create a new user, save user info into the database
// note that the password in User object is not encrypted, but in the database it is
func (user *User) Create() (err error) {
	// Postgres does not automatically return the last insert id, because it would be wrong to assume
	// you're always using a sequence.You need to use the RETURNING keyword in your insert to get this
	// information from postgres.
	statement := "insert into users (name, password, created_at) values ($1, $2, $3, $4) returning id"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	// use QueryRow to return a row and scan the returned id into the User struct
	// err = stmt.QueryRow(user.Name, Encrypt(user.Password), time.Now()).Scan(&user.Id, &user.Uuid, &user.CreatedAt)
	return
}

// get a single user given the name
func UserByName(name string) (user User, err error) {
	user = User{}
	err = Db.QueryRow("SELECT id, name, password, balance, batteryCapacity FROM users WHERE name = $1", name).
		Scan(&user.Id, &user.Name, &user.Password, &user.Balance, &user.BatteryCapacity)
	return
}

// TODO: legacy
func UserBySession(sess *Session) (user User, err error) {
	user = User{}
	err = Db.QueryRow("SELECT id, uuid, name, password, balance, batteryCapacity, created_at FROM users WHERE id = $1", sess.Uuid).
		Scan(&user.Id, &user.Name, &user.Password, &user.Balance, &user.BatteryCapacity)
	return
}

// Create a new session for an existing user
func (user *User) CreateSession() (session Session, err error) {
	statement := "insert into sessions (uuid, name, user_id, created_at) values ($1, $2, $3, $4) returning id, uuid, name, user_id, created_at"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	// use QueryRow to return a row and scan the returned id into the Session struct
	err = stmt.QueryRow(createUUID(), user.Name, user.Id, time.Now()).Scan(&session.Id, &session.Uuid, &session.Name, &session.UserId, &session.CreatedAt)
	return
}
