package controller

import (
	"encoding/json"
	"fmt"
	"hello/dbConfig"
	//"log"
	//"net/http"
)

func insertNodes() {
	json.Marshal(23)
	fmt.Println("Inside DB Controller")
	db := dbConfig.Connect()

	defer db.Close()

	//http.DefaultServeMux

}
