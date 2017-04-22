package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"sync"
)

var (
	connection            *sql.DB
	dbConnectionMutex     *sync.Mutex
)

func init() {
	dbConnectionMutex = &sync.Mutex{}
	connection = nil
}

func GetConnection() *sql.DB {
	dbConnectionMutex.Lock()
	if connection == nil {
		connection = connectToDb()
	}
	dbConnectionMutex.Unlock()
	return connection;
}

func connectToDb() *sql.DB {
	db, err := sql.Open("postgres", "user=discharm dbname=discharm password=J@spercat")
	if err != nil {
		log.Println("Unable to connect to database")
		log.Println(err)
		return nil
	}
	
	return db
}
