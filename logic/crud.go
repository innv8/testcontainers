package logic

import "database/sql"

type User struct {
	Id   int64
	Name string
	Age  int
}

func CreateUser(db *sql.DB, name string, age int) (id int64, err error) {
	var query = "INSERT INTO user (name, age) VALUES (?,?)"
	res, err := db.Exec(query, name, age)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func GetUser(db *sql.DB, id int64) (user User, err error) {
	var query = "SELECT id, name, age FROM user WHERE id=? LIMIT 1
	err = db.QueryRow(query, id).Scan(&user.Id, &user.Name, &user.Age)
	return user, err
}

func UpdateUser(db *sql.DB, id int64, name string, age int) (err error) {
	var query = "UPDATE user SET name=?, age=? WHERE id=? LIMIT 1"
	_, err = db.Exec(query, name, age, id)
	return err
}

func DeleteUser(db *sql.DB, id int64) (err error) {
	var query = "DELETE FROM user WHERE id=? LIMIT 1"
	_, err = db.Exec(query, id)
	return err
}
