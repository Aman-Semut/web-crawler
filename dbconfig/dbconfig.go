package dbconfig

import (
	"database/sql"
	"fmt"
	
	// _ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "aman"
	password = "python3"
	dbname   = "webcrawler"
  )

var db *sql.DB

func Connect() *sql.DB {
	dbDriver := "postgres"
	// dbUser := "aman"
	// dbPass := "python3"
	// dbName := "webcrawler"
	// dbIp := "127.0.0.1"
	// dbPort := 5432

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable",
    host, port, user, password, dbname)


	// db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp("+dbIp+":"+dbPort+")/"+dbName)
	db, err := sql.Open(dbDriver, psqlInfo)
	if err != nil {
		fmt.Println("Error connecting to database: ", err.Error())
		panic(err.Error())
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("Error connecting to database: ", err.Error())
		panic(err.Error())
	}

	fmt.Println("Successfully connected to database")



	return db
}

func main(){
	Connect();
	return;
}