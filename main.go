package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"

	"hello/graph"
	"hello/dbController"

	"github.com/steelx/extractlinks"
	// "strings"
	// "golang.org/x/net/html"
)

// add comments
// write more clean code add more white spaces
// delv debugger

type ReqData struct {
	Id  int64  `json:"id"`
	Url string `json:"url"`
}

type CustomError struct {
	err  error
	msg  string
	code int
}

// Declaring variables for use in crawler algorithm
var (
	//This is a channel of string having a maximum buffer size of 100, it stores the urls obtained from recursive calls and helps to add them to the url graph
	urlQueue = make(chan string, 100) // set an upper limit - Done

	//This is tls config to avoid certificate issues when visiting some of the urls
	config = &tls.Config{InsecureSkipVerify: true}

	//This is a simple http transport
	transport = &http.Transport{
		TLSClientConfig: config,
	}
	//Stores if the url is visited or not
	hasCrawled = make(map[string]bool)

	netClient *http.Client
	//A graph data structure to store the network visited in the for of graph and the urls as the nodes of the graph
	graphMap = graph.NewGraph() // when we move to DB check the usage of it
	// try to use noSQL DB
	// batch mode storing of size 100
	// find a optimal value for batch size

)

// Initialisation of http client and the go routine for handling interrupt signals
func init() {
	netClient = &http.Client{
		Transport: transport,
	}
	go SignalHandler(make(chan os.Signal, 1))
}

func handleReqs() {
	r := http.NewServeMux()

	r.HandleFunc("/post/", postCrawlRequest)
	r.HandleFunc("/get", getCrawlRequest)

	err := http.ListenAndServe(":9000", r)
	log.Fatal(err)
}

func main() {
	fmt.Println("Hello, playground")
	// dbController.AddDataToUser()
	handleReqs()
}

func getCrawlRequest(rw http.ResponseWriter, r *http.Request) {

	// fmt.Println("GET params were:", r.URL.Query())
	url := r.URL.Query().Get("url")
	fmt.Println("Url received:", url)

	getGraph := dbController.GetData(url)

	resp, err := json.Marshal(getGraph)
	//fmt.Println("Response : ", string(resp))
	if err != nil {
		// fmt.Println("Graphmap couldn't be responded", err)
		// http.Error(rw, err.Error(), 500)
		error := CustomError{err, err.Error(), 500}
		throwHttpError(rw, error)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	res, err := rw.Write(resp)

	if err != nil {
		// fmt.Println("Graphmap couldn't be responded duw to rw:", err)
		// http.Error(rw, err.Error(), 500)
		error := CustomError{err, err.Error(), 500}
		throwHttpError(rw, error)
		return
	}
	fmt.Println("Responding to postman with ", res)
}

func postCrawlRequest(rw http.ResponseWriter, req *http.Request) {

	b, err := io.ReadAll(req.Body)

	defer req.Body.Close()
	if err != nil {
		http.Error(rw, err.Error(), 500)
		return
	}

	// Unmarshal
	var body ReqData
	err = json.Unmarshal(b, &body)
	if err != nil {
		http.Error(rw, err.Error(), 500)

		return
	}

	//If there are no arguments , just make  a safe exit
	if len(body.Url) < 1 {
		fmt.Println("URL is missing")
		// os.Exit(1)
		//safeExit(1)
		var myError CustomError
		myError.msg = "URL is missing"
		myError.code = 400
		throwHttpError(rw, myError)
		return
	}
	// check with len args < 1 - Done

	//get the base url from the slice of urls i.e args
	baseURL := body.Url // rename the baseURL - Done

	//Append the baseURL to the urlQueue
	go func() {
		urlQueue <- baseURL
	}()

	// concurrency handling -TODO
	var maxdepth int = 2
	for href := range urlQueue {
		if !hasCrawled[href] {
			// this needs to be a go routine and set upper limit on the number of go routines-TODO

			crawlLink(href, rw)
			maxdepth--
			if maxdepth == 0 {
				resp, err := json.Marshal(graphMap.Adjacency)
				//fmt.Println("Response : ", string(resp))
				if err != nil {
					// fmt.Println("Graphmap couldn't be responded", err)
					// http.Error(rw, err.Error(), 500)
					error := CustomError{err, err.Error(), 500}
					throwHttpError(rw, error)
					return
				}

				rw.WriteHeader(http.StatusOK)
				rw.Header().Set("Content-Type", "application/json; charset=utf-8")
				res, err := rw.Write(resp)

				if err != nil {
					// fmt.Println("Graphmap couldn't be responded duw to rw:", err)
					// http.Error(rw, err.Error(), 500)
					error := CustomError{err, err.Error(), 500}
					throwHttpError(rw, error)
					return
				}
				fmt.Println("Responding to postman with ", res)

				

				//  add data 
				dbController.AddData(graphMap.Adjacency, baseURL)

				return
				// safeExit(0)
			}

		}
	}
}

// This function is called when some signal is received from the os
func SignalHandler(c chan os.Signal) {
	signal.Notify(c, os.Interrupt)
	// while loop would be more clean here ## check
	s := <-c

	// for s := <-c; ; s = <-c {

	//Based on the type of interrupt signal received , do something

	switch s {
	case os.Interrupt:
		fmt.Println("^C received")
		fmt.Println("<----------- ----------- ----------- ----------->")
		fmt.Println("<----------- ----------- ----------- ----------->")
		//If ^C received , complete the path between src and dest and make a safe exit
		graphMap.CreatePath("https://youtube.com/", "https://youtube.com/YouTubeRedOriginals")
		// safe exit -> check for defer function # db connection - Partially Done
		//os.Exit(0)
		safeExit(0)

	case os.Kill:
		fmt.Println("SIGKILL received")
		// os.Exit(1)
		safeExit(1)
	}
	// }
}

// This function crawls from a given baseURL and compltes the network map after crwaling
func crawlLink(baseHref string, rw http.ResponseWriter) {
	//Adds the current node or baseHref to the network map
	graphMap.AddVertex(baseHref) // check this for exit
	hasCrawled[baseHref] = true

	//Starts the crawling process from the given baseHref
	fmt.Println("Crawling... ", baseHref) // use logger - TODO

	//Gets the response from the netClient,checks for errors and processes the responses
	resp, err := netClient.Get(baseHref)
	checkErr(err, rw)

	//close each network connection after crawling the network from a strating baseHref
	defer resp.Body.Close()

	//Extract the links from the reposnse body which contains HTML
	links, err := extractlinks.All(resp.Body)
	checkErr(err, rw)

	// get all texts and use the information - TODO
	// texts := buildText(resp.Body)
	// fmt.Println(texts)
	// os.Exit(0)

	//Loop over the extracted links from the reposnse body
	for _, l := range links {
		//If link is not found or is empty , do nothing
		if l.Href == "" {
			continue
		}

		//Get the fixedURL from the link
		fixedUrl := toFixedUrl(l.Href, baseHref)

		//Only If the link is obtained from the original baseHref i.e the link is from the same domain as baseHref, add it to the graph
		if baseHref != fixedUrl {
			graphMap.AddEdge(baseHref, fixedUrl)
		}

		// try using this -> urlQueue <- url

		go func(url string) {
			urlQueue <- fixedUrl
		}(fixedUrl)
	}
}

func checkErr(err error, rw http.ResponseWriter) {
	if err != nil {
		fmt.Println("Some error occured : ", err)
		// safe exit -> continue
		// os.Exit(0)
		myError := CustomError{err, err.Error(), 400}
		throwHttpError(rw, myError)
		//safeExit(0)
		return

	}
}

// Get the fixedUrl from the href
func toFixedUrl(href, base string) string {
	uri, err := url.Parse(href)

	//If there is some error or the url is mail url or telphone url, just return the origin base
	if err != nil || uri.Scheme == "mailto" || uri.Scheme == "tel" {
		return base
	}

	baseURL, err := url.Parse(base)
	if err != nil {
		return ""
	}

	uri = baseURL.ResolveReference(uri)
	return uri.String()
}

// Write Code for safeExit store in DB, close connections etc - TODO
func safeExit(code int) string {

	defer os.Exit(code)

	//fmt.Println("Printing graph before exiting:")
	//graphMap.Print()

	fmt.Println("Closing connections...")

	switch code {
	case 0:
		// insert data
		//Do something
		break
	case 1:
		//Do something
		break
	default:
		//Do something
		break
	}

	fmt.Println("Safe Exit with code " + strconv.Itoa(code))

	return "Safe Exit with code " + strconv.Itoa(code)
}

func throwHttpError(rw http.ResponseWriter, customError CustomError) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(customError.code)
	json.NewEncoder(rw).Encode(map[string]string{"error": customError.msg})

}
