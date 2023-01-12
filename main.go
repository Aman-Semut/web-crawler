package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"

	"hello/graph"

	"github.com/steelx/extractlinks"
	// "strings"
	// "golang.org/x/net/html"
	// "io"
)

// add comments
// write more clean code add more white spaces
// delv debugger

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

func main() {
	// package cobra for command line arguments
	args := os.Args[1:]
	
	//If there are no arguments , just make  a safe exit
	if len(args) < 1 {
		fmt.Println("URL is missing")
		// os.Exit(1)
		safeExit(1)
	}
	// check with len args < 1 - Done

	//get the base url from the slice of urls i.e args
	baseURL := args[0] // rename the baseURL - Done

	//Append the baseURL to the urlQueue
	go func() {
		urlQueue <- baseURL
	}()

	// concurrency handling -TODO
	for href := range urlQueue {
		if !hasCrawled[href] {
			// this needs to be a go routine and set upper limit on the number of go routines-TODO
			crawlLink(href)
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
func crawlLink(baseHref string) {
	//Adds the current node or baseHref to the network map
	graphMap.AddVertex(baseHref) // check this for exit
	hasCrawled[baseHref] = true

	//Starts the crawling process from the given baseHref
	fmt.Println("Crawling... ", baseHref) // use logger - TODO

	//Gets the response from the netClient,checks for errors and processes the responses
	resp, err := netClient.Get(baseHref)
	checkErr(err)

	//close each network connection after crawling the network from a strating baseHref
	defer resp.Body.Close()

	//Extract the links from the reposnse body which contains HTML
	links, err := extractlinks.All(resp.Body)
	checkErr(err)

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
			urlQueue <- url
		}(fixedUrl)
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		// safe exit -> continue
		// os.Exit(0)
		safeExit(0)
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

// func buildText(htmlBody io.Reader) string {
// 	n, err := html.Parse(htmlBody)
// 	if err != nil {
// 		return ""
// 	}

// 	fmt.Println(n.Data)
// 	// os.Exit(0)
// 	return ""
// 	// if n.Type == html.TextNode {
// 	// 	return n.Data
// 	// }

// 	// var text string
// 	// for c := n.FirstChild; c != nil; c = c.NextSibling {
// 	// 	text += buildText(c)
// 	// }
// 	// return strings.Join(strings.Fields(text), " ")
// }

// Write Code for safeExit store in DB, close connections etc - TODO
func safeExit(code int) string {

	defer os.Exit(code)

	fmt.Println("Printing graph before exiting:")
	graphMap.Print()

	fmt.Println("Closing connections...")

	switch code {
	case 0:

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
