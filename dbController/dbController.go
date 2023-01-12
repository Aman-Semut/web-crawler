package dbController

import (
	"encoding/json"
	"fmt"
	"hello/dbconfig"
	// "hello/graph"
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

func AddData(graph map[string][]string) {
	db := dbconfig.Connect()

	fmt.Println("Inside DB Controller")

	length := len(graph)
	for key, value := range graph {
		fmt.Println("Key: ", key, "Value: ", value)
		//fmt.Println("Value: ", value)
		//fmt.Println("Length: ", length)
		length--
		if length == 0 {
			break
		}
	}



	defer db.Close()


}
