package dbController

import (
	// "encoding/json"
	"fmt"
	"hello/dbconfig"
	// "hello/graph"
	//"log"
	//"net/http"
)

func AddData(graph map[string][]string, baseURL string) bool {
	db := dbconfig.Connect()

	fmt.Println("Inside DB Controller")

	length := len(graph)
	id := 1
	// columns := ("id", "base_url", "current_url", "parent_url")

	values := make([]interface{}, 0, length*3)

	done := true

	for key, value := range graph {
		// fmt.Println("Key: ", key, "Value: ", value)

		for _, v := range value {
			values = append(values, baseURL, v, key)
			id++
		}

		length--
		if length == 0 {
			break
		}
	}

	query_string := `INSERT INTO stored_links (base_url, current_url, parent_url) VALUES ($1, $2, $3)`

	for i := 0; i < len(values); i += 3 {
		fmt.Println("Values: ", values[i], values[i+1], values[i+2])
		_, e := db.Exec(query_string, values[i], values[i+1], values[i+2])
		if e != nil {
			fmt.Println("Error executing statement: ", e.Error())
			done = false
		}
	}

	defer db.Close()
	return done
}

func GetData(baseURL string) map[string][]string {
	db := dbconfig.Connect()

	defer db.Close()

	query_string := `SELECT * FROM stored_links WHERE base_url = '` + baseURL + `'`

	rows, err := db.Query(query_string)
	if err != nil {
		fmt.Println("Error executing statement: ", err.Error())
	}

	defer rows.Close()

	graph := make(map[string][]string)

	for rows.Next() {
		var id int
		var base_url string
		var current_url string
		var parent_url string

		err := rows.Scan(&id, &base_url, &current_url, &parent_url)
		if err != nil {
			fmt.Println("Error executing statement: ", err.Error())
		}

		graph[parent_url] = append(graph[parent_url], current_url)
	}

	fmt.Println("Graph: ", graph)

	return graph

}
