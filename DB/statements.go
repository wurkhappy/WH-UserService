package DB

import (
	"database/sql"
	_ "github.com/bmizerany/pq"
	// "log"
)

var SaveUser *sql.Stmt
var UpdateUser *sql.Stmt
var UpsertUser *sql.Stmt
var FindUserByEmail *sql.Stmt
var FindUserByID *sql.Stmt
var DeleteUser *sql.Stmt
var FindUsers *sql.Stmt

func CreateStatements() {
	var err error
	SaveUser, err = DB.Prepare("INSERT INTO wh_user(id, password, data) VALUES($1, $2, $3)")
	if err != nil {
		panic(err)
	}

	UpdateUser, err = DB.Prepare("UPDATE wh_user SET password = $2, data = $3 WHERE id = $1")
	if err != nil {
		panic(err)
	}

	UpsertUser, err = DB.Prepare("SELECT upsert_user($1, $2, $3)")
	if err != nil {
		panic(err)
	}

	FindUserByEmail, err = DB.Prepare("SELECT password, data FROM wh_user WHERE data->>'email' = $1")
	if err != nil {
		panic(err)
	}

	FindUserByID, err = DB.Prepare("SELECT password, data FROM wh_user WHERE id = $1")
	if err != nil {
		panic(err)
	}

	DeleteUser, err = DB.Prepare("DELETE FROM wh_user WHERE id = $1")
	if err != nil {
		panic(err)
	}

	FindUsers, err = DB.Prepare("SELECT data FROM wh_user WHERE id = ANY($1)")
	if err != nil {
		panic(err)
	}
}
