package dbconfig

import (
	"database/sql"
	"fmt"
	
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func Connect() *sql.DB {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "pyThon@3"
	dbName := "webcrawler"
	dbIp := "127.0.0.1"
	dbPort := "3306"

	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp("+dbIp+":"+dbPort+")/"+dbName)
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