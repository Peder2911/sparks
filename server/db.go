/*
SPRAWL uses an SQLite database. This is fine for now.
*/
package main

import (
   "database/sql"
   "fmt"
   _ "github.com/mattn/go-sqlite3"
)

func db_connect() (*sql.DB) {
   db, err := sql.Open("sqlite3","./database.db")
   if err != nil {
      panic(fmt.Sprintf("Failed to connect to DB: %s", err))
   }
   return db
}

func initialize(connect func()(*sql.DB)) {
   db := connect()
   defer db.Close()
   db.Exec(`
   create table users (
      id integer not null primary key autoincrement,
      name text unique,
      password text)
   `)
   db.Exec(`
   create table sessions (
      id text not null primary key,
      user int,
      foreign key (user) references user(id)
   )
   `)
}
