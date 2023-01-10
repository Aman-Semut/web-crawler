package main

import (
	"crypto/tls"
	"fmt"
	"github.com/steelx/extractlinks"
	"github.com/steelx/webscrapper/graph"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	// "strings"
	// "golang.org/x/net/html"
	// "io"
)

// add comments
// write more clean code add more white spaces
// delv debugger

var (
	urlQueue  = make(chan string) // set an upper limit
	config    = &tls.Config{InsecureSkipVerify: true}
	transport = &http.Transport{
		TLSClientConfig: config,
	}
	hasCrawled = make(map[string]bool) // difference between array and slice 
	netClient  *http.Client
	graphMap   = graph.NewGraph() // when we move to DB check the usage of it
	// try to use noSQL DB
	// batch mode storing of size 100 
	// find a optimal value for batch size
)

func init() {
	netClient = &http.Client{
		Transport: transport,
	}
	go SignalHandler(make(chan os.Signal, 1))
}

func main() {
	// package cobra for command line arguments
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("URL is missing")
		os.Exit(1)
	}
	// check with len args > 1

	baseUrl := args[0] // rename the baseURL

	go func() {
		urlQueue <- baseUrl
	}()

	// concurrency handling 
	for href := range urlQueue {
		if !hasCrawled[href] {
			// this needs to be a go routine and set upper limit on the number of go routines
			crawlLink(href)
		}
	}

}

func SignalHandler(c chan os.Signal) {
	signal.Notify(c, os.Interrupt)
	// while loop would be more clean here ## check
	s := <-c

	// for s := <-c; ; s = <-c {
		switch s {
		case os.Interrupt:
			fmt.Println("^C received")
			fmt.Println("<----------- ----------- ----------- ----------->")
			fmt.Println("<----------- ----------- ----------- ----------->")
			graphMap.CreatePath("https://youtube.com/", "https://youtube.com/YouTubeRedOriginals")
			// safe exit -> check for defer function # db connection 
			os.Exit(0)
		case os.Kill:
			fmt.Println("SIGKILL received")
			os.Exit(1)
		}
	// }
}

func crawlLink(baseHref string) {
	graphMap.AddVertex(baseHref) // check this for exit
	hasCrawled[baseHref] = true
	fmt.Println("Crawling... ", baseHref) // use logger 
	resp, err := netClient.Get(baseHref)
	checkErr(err)
	defer resp.Body.Close()

	links, err := extractlinks.All(resp.Body)
	checkErr(err)

	// get all texts and use the information
	// texts := buildText(resp.Body)
	// fmt.Println(texts)
	// os.Exit(0)

	for _, l := range links {
		if l.Href == "" {
			continue
		}
		fixedUrl := toFixedUrl(l.Href, baseHref)
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
		os.Exit(0)
	}
}

func toFixedUrl(href, base string) string {
	uri, err := url.Parse(href)
	if err != nil || uri.Scheme == "mailto" || uri.Scheme == "tel" {
		return base
	}

	baseUrl, err := url.Parse(base)
	if err != nil {
		return ""
	}

	uri = baseUrl.ResolveReference(uri)
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

