package dbController

import (
	"encoding/json"
	"fmt"
	"hello/dbconfig"
)

func AddDataToUser() {
	db := dbconfig.Connect()

	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO user (username) VALUES (?)")
	if err != nil {
		fmt.Println("Error preparing statement: ", err.Error())
		panic(err.Error())
	}

	res, err := stmt.Exec("aman")
	if err != nil {
		fmt.Println("Error executing statement: ", err.Error())
		panic(err.Error())
	}
	return nil
}
