package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"github.com/spf13/viper"
	"sync"
	"fmt"
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
	user := viper.GetString("DatabaseUsername")
	password := viper.GetString("DatabasePassword")
	name := viper.GetString("DatabaseName")
	host := viper.GetString("DatabaseHost")
	
	db, err := sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s password=%s host=%s", user, name, password, host))
	if err != nil {
		log.Println("Unable to connect to database")
		log.Println(err)
		return nil
	}
	
	return db
}
