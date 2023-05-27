package data

type User struct {
	Id              int
	Uuid            string
	Name            string
	Password        string
	IsAdmin         bool // unused, Id == 0 -> Admin
	Balance         float64
	BatteryCapacity float64
}

// create a new user, save user info into the database
// note that the password in User object is not encrypted, but in the database it is
func (user *User) Create(isAdmin bool) (err error) {
	// Postgres does not automatically return the last insert id, because it would be wrong to assume
	// you're always using a sequence.You need to use the RETURNING keyword in your insert to get this
	// information from postgres.
	statement := "insert into users (uuid, name, password, isadmin, balance, batteryCapacity) values ($1, $2, $3, $4, $5, $6) returning id"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	// use QueryRow to return a row and scan the returned id into the User struct
	err = stmt.QueryRow(createUUID(), user.Name, Encrypt(user.Password), user.IsAdmin, user.Balance, user.BatteryCapacity).Scan(&user.Id)
	return
}

// get a single user given the name
func UserByName(name string) (user User, err error) {
	user = User{}
	err = Db.QueryRow("SELECT * FROM users WHERE name = $1", name).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Password, &user.IsAdmin, &user.Balance, &user.BatteryCapacity)
	return
}

// get a single user given the id
func UserById(id int) (user User, err error) {
	user = User{}
	err = Db.QueryRow("SELECT * FROM users WHERE id = $1", id).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Password, &user.IsAdmin, &user.Balance, &user.BatteryCapacity)
	return
}
