package data

import "time"

type Session struct {
	Id        int
	Uuid      string
	Name      string
	UserId    int
	CreatedAt time.Time
}

// check if session is valid in the database
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

func UserBySession(sess *Session) (user User, err error) {
	user = User{}
	err = Db.QueryRow("SELECT * FROM users WHERE id = $1", sess.Uuid).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Password, &user.Balance, &user.BatteryCapacity)
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