package main

import (
	"bank/db"
	"bank/db/memdb"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	memdb.Connect()
	newdb := memdb.Newdb
	stmt, err := newdb.Prepare("create table if not exists users (Id INTEGER PRIMARY KEY AUTOINCREMENT , Name text not null , Email text);")
	if err != nil {
		panic(err)
	}
	stmt.Exec()
	http.HandleFunc("/users", userHandler)
	http.ListenAndServe(":8080", nil)
}

func userHandler(rw http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		// TODO: crete user
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error())) //
			return
		}
		u := &db.User{}
		err = json.Unmarshal(data, u)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(err.Error()))
			return
		}
		u, err = db.CreateUser(r.Context(), u.Name, u.Email)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}
		data, err = json.Marshal(u)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}
		rw.Write(data)
	}

	// get all users
	if r.Method == http.MethodGet {
		users, err := db.ListUsers(r.Context())
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}
		data, err := json.Marshal(users)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}
		rw.Write(data)
	}

	// delete user by id
	if r.Method == http.MethodDelete {
		id := r.URL.Query().Get("ID")
		newid, _ := strconv.ParseInt(id, 0, 0)
		err := db.DeleteUser(r.Context(), newid)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}
		rw.Write([]byte("user is deleted"))
	}

	if r.Method == http.MethodPatch {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}
		u := &db.User{}
		err = json.Unmarshal(data, u)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(err.Error()))
			return
		}
		// user by id
		us, err := db.UserById(r.Context(), u.Id)
		if errors.Is(err, db.ErrNotFound) {
			rw.WriteHeader(http.StatusNotFound)
			rw.Write([]byte(err.Error()))
			return
		}
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}
		changed := false
		if u.Name != us.Name {
			changed = true
		}
		if u.Email != us.Email {
			changed = true
		}
		if !changed {
			return
		}

		// update user
		_, err = db.UpdateUser(r.Context(), u.Id, u.Name, u.Email)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}
		data, err = json.Marshal(u)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}
		rw.Write(data)
	}
}
