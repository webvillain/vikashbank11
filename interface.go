package db

import (
	"bank/db/memdb"
	"context"
	"errors"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var ErrNotFound = errors.New("not found")

var ErrNotImplemented = errors.New("not implemnted")

func UserById(ctx context.Context, id int64) (*User, error) {
	var getuserbyid *User
	getuserbyidquery := `SELECT * FROM users WHERE ID = ?`
	row, err := memdb.Newdb.Query(getuserbyidquery, id)
	if err != nil {
		panic(err)
	}
	if row.Next() {
		row.Scan(&getuserbyid.Id, &getuserbyid.Name, &getuserbyid.Email)

	}
	return getuserbyid, nil
}

func CreateUser(ctx context.Context, name string, email string) (*User, error) {
	var newuser *User
	newuserstmt := `INSERT INTO users (Name , Email)VALUES(?,?)`
	stmt, err := memdb.Newdb.Prepare(newuserstmt)
	if err != nil {
		panic(err)
	}
	res, err := stmt.Exec(name, email)
	if err != nil {
		log.Fatal(err)
	}
	id, _ := res.LastInsertId()
	// /rowsaffected, _ := res.RowsAffected()
	newuser.Id = int64(id)
	newuser.Name = name
	newuser.Email = email
	return newuser, nil
}

func UpdateUser(ctx context.Context, id int64, name string, email string) (*User, error) {

	var updateduserbyidquery = `UPDATE users
	SET Name = ?,
		Email = ?
	WHERE
		Id = ? ;`
	var updateduser *User
	stmt, err := memdb.Newdb.Prepare(updateduserbyidquery)
	if err != nil {
		panic(err)
	}
	res, err := stmt.Exec(name, email, id)
	if err != nil {
		panic(err)
	}
	fmt.Println(res.RowsAffected())

	return updateduser, nil
}

func DeleteUser(ctx context.Context, id int64) error {
	deleteuserbyidquery := `DELETE FROM users WHERE ID = ?`
	stmt, err := memdb.Newdb.Prepare(deleteuserbyidquery)
	if err != nil {
		panic(err)
	}
	res, err := stmt.Exec(id)
	if err != nil {
		panic(err)
	}
	fmt.Println(res.RowsAffected())
	return nil
}

func ListUsers(ctx context.Context) ([]*User, error) {
	var user *User
	var users []*User
	rows, err := memdb.Newdb.Query("SELECT * FROM users;")
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		rows.Scan(&user.Id, &user.Name, &user.Email)
	}

	users = append(users, user)
	fmt.Println("Total No. Of Users :", len(users))

	return users, nil
}
