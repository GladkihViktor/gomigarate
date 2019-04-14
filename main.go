package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	os "os"

	_ "github.com/lib/pq"
)

// Config for database
type Config struct {
	Driver   string
	Host     string
	Port     string
	User     string
	Dbname   string
	Password string
	Sslmode  string
}

//Migrations realtaion table Migrations
type Migrations struct {
	id       int
	filename string
	applied  bool
}

var (
	//Db pointer to sql.DB
	Db *sql.DB
	//StrConnect string for connection Db
	StrConnect string
	//DbConfig pointer to struct Config
	DbConfig *Config
)

// GetConfigDB ...
func GetConfigDB() *Config {
	var config Config

	conf, errConf := os.Open("config.json")
	if errConf != nil {
		panic(errConf)
	}
	defer conf.Close()

	decode := json.NewDecoder(conf)
	err := decode.Decode(&config)
	if err != nil {
		panic(err)
	}
	return &config
}

//GetStrConnect ....
func GetStrConnect(c Config) string {
	strConnect := "host=" + c.Host + " port=" + c.Port + " user=" + c.User + " dbname=" + c.Dbname + " password=" + c.Password + " sslmode=" + c.Sslmode
	return strConnect
}

func init() {
	DbConfig = GetConfigDB()
	StrConnect = GetStrConnect(*DbConfig)
}

func initDb() *sql.DB {
	Db, err := sql.Open(DbConfig.Driver, StrConnect)
	if err != nil {
		panic(err)
	}
	return Db
}

//ReadFileSQL ...
func ReadFileSQL(path string) string {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(file)
}

func main() {
	fmt.Println("\x1b[1;34m", "[START MIGRATE] ", "\x1b[0m")
	files, err := ioutil.ReadDir("source")
	if err != nil {
		log.Fatalln("\x1b[1;31m", "[ERROR] ", "\x1b[0m", err)
	}
	db := initDb()
	tx, err := db.Begin()

	for _, file := range files {
		if file.IsDir() {
			fmt.Println("\x1b[1;34m", "[END]", "\x1b[0m")
			return
		}
		script := ReadFileSQL("source/" + file.Name())

		rows, err := db.Query("select check_migrate($1)", file.Name())
		defer rows.Close()
		if err != nil {
			fmt.Println("\x1b[1;31m", "[ERROR] ", "\x1b[0m", err)
		}

		var idMigrate int
		rows.Scan(&idMigrate)

		if idMigrate != 0 {
			fmt.Println("\x1b[1;34m", "[START] ", "\x1b[0m", file.Name())
			_, err = tx.Exec(script)
			if err != nil {
				tx.Rollback()
				fmt.Println("\x1b[1;31m", "[ERROR] ", "\x1b[0m", err)
			}

			fmt.Println("\x1b[1;32m", "[OK] ", "\x1b[0m", file.Name())

			str := file.Name()
			fmt.Println(str)

			res, err := db.Query("select insert_migration($1)", file.Name())
			defer res.Close()

			if err != nil {
				fmt.Println("\x1b[1;31m", "[ERROR] ", "\x1b[0m", err)
			}

			var newID int
			res.Scan(&newID)

			fmt.Println("PK for row migration: ", newID)
		}
	}

	tx.Commit()
	fmt.Println("\x1b[1;34m", "[END]", "\x1b[0m")
}
