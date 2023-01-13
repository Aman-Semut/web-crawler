package dbController

import (
	// "encoding/json"
	"fmt"
	"hello/dbconfig"
)

// func AddDataToUser() {
// 	db := dbconfig.Connect()

// 	defer db.Close()

// 	stmt, err := db.Prepare("INSERT INTO user (id, username) VALUES (?, ?)")
// 	if err != nil {
// 		fmt.Println("Error preparing statement: ", err.Error())
// 		panic(err.Error())
// 	}

// 	res, err := stmt.Exec(1, 'a')
// 	if err != nil {
// 		fmt.Println("Error executing statement: ", err.Error())
// 		panic(err.Error())
// 	}
	
// 	fmt.Println(res)
// }


func AddDataToUser() {
	db := dbconfig.Connect()

	defer db.Close()

	link := "https://www.google.com"
	query_string := `INSERT INTO links (link) VALUES ('` + link + `')`

	_, e := db.Exec(query_string)
	if e != nil {
		fmt.Println("Error executing statement: ", e.Error())
	}


}
