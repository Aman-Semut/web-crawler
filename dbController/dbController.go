package dbController

import (
	"encoding/json"
	"fmt"
	"hello/dbconfig"
	//"log"
	//"net/http"
)

func insertNodes() {
	json.Marshal(23)
	fmt.Println("Inside DB Controller")
	db := dbconfig.Connect()

	defer db.Close()

	//http.DefaultServeMux

}


// add data 
// get data 

func AddData() {
	db := dbconfig.Connect()

	defer db.Close()


}
