package utils

import (
	"fmt"
	"io"
	"os"
	"sync"
)

func AsyncDownloadFile(path string, fileName string, url string, wg *sync.WaitGroup, proxy bool) {
	defer wg.Done()

	if url == "" {
		return
	}
	
	file, err := os.Open(path+fileName)
	stat, err := file.Stat()
	if err != nil {
         fmt.Printf("file read error %v for %v", err.Error(), fileName)
     }
	if stat.Size() > 0 {
		file.Close()
		return
	}
	file.Close()
	
	// Create the file
	out, err := os.Create(path + fileName)
	if err != nil {
		fmt.Println(path + fileName)
		return
	}
	defer out.Close()

	// Get the data
	hc := new(HttpConnection)
	hc.Init(proxy)
	res, err := hc.GetHttpResponse(url)
	if err != nil {
		fmt.Println("http error on " + url)
		return
	}
	defer res.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, res.Body)
	if err != nil {
		fmt.Println("copy fail")
		return
	}
	fmt.Println("wg done")
	return
}
