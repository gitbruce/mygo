package main

import (
	"bruce/utils"
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sync"
	"strings"
)

const PdfPath = "/tmp/patent/"
const PatentCsv = "patents.csv"

type PatentDetail struct {
	Prefix            string
	ApplicationNumber string
	Title             string
	AssigneeOriginal  string
	Language          string
	PriorityDate      string
	FilingDate        string
	PublicationDate   string
	GrantDate         string
	Inventor          string
	Assignee          string
	Pdf               string
}

func main() {
	csvFile, _ := os.Open(PdfPath + PatentCsv)
	reader := csv.NewReader(bufio.NewReader(csvFile))
	reader.LazyQuotes = true
	var patent PatentDetail
	var i = 0
	var wg sync.WaitGroup
	var fileName string
	for ; ; i++ {
		line, error := reader.Read()
		if error == io.EOF {
			fmt.Println("break")
			break
		} else if error == csv.ErrFieldCount {
			fmt.Println("nmumber file is different")
		} else if error != nil {
			fmt.Println("fatal")
			//			log.Fatal(error.Error())
		}
		if i < 1 {
			continue
		}

		patent = PatentDetail{
			Prefix:            line[0],
			ApplicationNumber: line[1],
			Language:          line[4],
			Pdf:               line[11],
		}
		url := strings.TrimSpace(patent.Pdf)
		if (url == "") {
			continue
		}
		fileName = patent.Prefix + "_" + strings.TrimSpace(patent.ApplicationNumber) + ".pdf"
		wg.Add(1)
		fmt.Println("downloading " + url +" to " + fileName)
		go utils.AsyncDownloadFile(PdfPath, fileName, url, &wg, true)
	}
	wg.Wait()
}
