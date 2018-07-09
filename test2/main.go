package main

import (
	"sync"
	"net/http"
	"fmt"
)

func checkUrl(url string,  wg *sync.WaitGroup) {
	defer wg.Done()
			// Fetch the URL.
			http.Get(url)
		fmt.Println("done")
}

func main() {
	var wg sync.WaitGroup
	var urls = []string{
		"http://www.baidu.com/",
		"http://www.baidu.com/",
		"http://www.baidu.com/",
	}
	for _, url := range urls {
		// Increment the WaitGroup counter.
		fmt.Println("adding")
		wg.Add(1)
		// Launch a goroutine to fetch the URL.
		checkUrl(url, &wg)
	}
	// Wait for all HTTP fetches to complete.
		fmt.Println("waiting")
	wg.Wait()
		fmt.Println("finish")
}
