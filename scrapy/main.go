package main

import (
	"bruce/patent"
	"bruce/utils"
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/tidwall/gjson"
	"log"
	"os"
	"strings"
	//	"sync"
)

const PdfPath = "/tmp/patent/"
const PatentCsv = "patents.csv"

var proxy bool

func init() {
	log.Println("Starting ...")
	flag.BoolVar(&proxy, "proxy", true, "whether to use proxy")
	flag.Parse()
	if proxy {
		log.Println("using proxy ")
	}
}

func replaceComma(val string, ns string) string {
	return strings.Replace(val, ",", ns, -1)
}

func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}

	return (err != nil)
}

func deleteFile(path string) {
	// delete file
	var err = os.Remove(path)
	if isError(err) {
		return
	}

	fmt.Println("==> done deleting file")
}

func getPatentTree(searchUrl string, sn int) int {
	f, err := os.OpenFile(PdfPath+PatentCsv, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	w := bufio.NewWriter(f)
	var fileError error
	total := 0
	results := getPatentResults(searchUrl)
	//	var wg sync.WaitGroup
	//	var fileName string
	if (sn == 1) {
		_, fileError = w.WriteString(patent.PatentHeader() + "\n")
		isError(fileError)
	}

	for i, result := range results {
		patentNumber := result.Publication_number
		mainPatent := getPatentDetail(patentNumber)
		mainPatent.Prefix = utils.PrefixString(float64(sn) + float64(i))
		_, fileError = w.WriteString(mainPatent.String() + "\n")
		isError(fileError)
		//		fileName = mainPatent.Prefix + "_" + mainPatent.ApplicationNumber + ".pdf"
		//		wg.Add(1)
		//		go utils.AsyncDownloadFile(PdfPath, fileName, mainPatent.Pdf, &wg)
		total++
		fmt.Println(mainPatent)
		subPatents := mainPatent.RelevantPatents
		for j, subPatent := range subPatents {
			subPatent := getPatentDetail(subPatent.ApplicationNumber)
			subPatent.Prefix = utils.PrefixString(float64(sn) + float64(i) + (float64(j)+1)/10.0)
			_, fileError = w.WriteString(subPatent.String() + "\n")
			isError(fileError)
			//			fileName = subPatent.Prefix + "_" + subPatent.ApplicationNumber + ".pdf"
			//			wg.Add(1)
			//			go utils.AsyncDownloadFile(PdfPath, fileName, subPatent.Pdf, &wg)
		}
		w.Flush()
		//		wg.Wait()
	}
	if fileError != nil {
		fmt.Println("got error")
	}
	return total
}

func getPatentResults(searchUrl string) []patent.PatentResults {
	googlePatents := utils.HttpJson(searchUrl, proxy)
	results := gjson.Get(googlePatents, "results.cluster.0.result.#.patent")
	patents := make([]patent.PatentResults, 0)
	json.Unmarshal([]byte(results.String()), &patents)
	return patents
}

func getPatentDetail(patentNumber string) patent.PatentDetail {
	hc := new(utils.HttpConnection)
	hc.Init(proxy)
	var patentUrl = "https://patents.google.com/patent/" + patentNumber + "/"
	fmt.Println("navigating " + patentUrl)
	res, err := hc.GetHttpResponse(patentUrl)
	if err != nil {
		fmt.Println("error to retieve url " + patentUrl)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	var mainPatent = new(patent.PatentDetail)
	mainPatent.ApplicationNumber = patentNumber
	mainPatent.Title = doc.Find("article h1+span").Text()
	mainPatent.Title = strings.TrimSpace(mainPatent.Title)
	mainPatent.Title = strings.TrimSuffix(mainPatent.Title, "\n")
	mainPatent.Title = replaceComma(mainPatent.Title, " ")
	link, _ := doc.Find("article a[itemprop=pdfLink]").Attr("href")
	mainPatent.Pdf = link
	var authors string
	doc.Find("article dd[itemprop=inventor]").Each(func(i int, selection *goquery.Selection) {
		if i == 0 {
			authors += selection.Text()
		} else {
			authors += "-" + selection.Text()

		}
	})
	mainPatent.Inventor = replaceComma(authors, "-")
	mainPatent.AssigneeOriginal = doc.Find("article dd[itemprop=assigneeOriginal]").Text()
	mainPatent.AssigneeOriginal = replaceComma(mainPatent.AssigneeOriginal, "")
	mainPatent.PriorityDate = doc.Find("article time[itemprop=priorityDate]").Text()
	mainPatent.FilingDate = doc.Find("article time[itemprop=filingDate]").First().Text()
	mainPatent.PublicationDate = doc.Find("article time[itemprop=publicationDate]").First().Text()
	mainPatent.GrantDate = doc.Find("article time[itemprop=grantDate]").First().Text()
	mainPatent.RelevantPatents = []patent.PatentDetail{}
	doc.Find("article tr[itemprop=appsClaimingPriority]").Each(func(i int, selection *goquery.Selection) {
		var relevantPatent = patent.PatentDetail{}
		var appNumber = selection.Find("span[itemprop=representativePublication]").Text()
		if appNumber != patentNumber {
			relevantPatent.ApplicationNumber = appNumber
			relevantPatent.Language = selection.Find("span[itemprop=primaryLanguage]").Text()
			relevantPatent.PriorityDate = selection.Find("td[itemprop=priorityDate]").Text()
			relevantPatent.FilingDate = selection.Find("td[itemprop=filingDate]").Text()
			relevantPatent.Title = selection.Find("td[itemprop=title]").Text()
			relevantPatent.Title = strings.TrimSpace(relevantPatent.Title)
			relevantPatent.Title = strings.TrimSuffix(relevantPatent.Title, "\n")
			relevantPatent.Title = replaceComma(relevantPatent.Title, "")
			link2, _ := selection.Find("a").Attr("href")
			relevantPatent.Pdf = link2
			mainPatent.RelevantPatents = append(mainPatent.RelevantPatents, relevantPatent)
		}
	})
	return *mainPatent
}

func main() {
	deleteFile(PdfPath + PatentCsv)
	var urls = []string{
		"https://patents.google.com/xhr/query?url=inventor%3D%E5%88%98%E7%91%BE%26assignee%3D%E4%B8%8A%E6%B5%B7%E8%B4%9D%E5%B0%94%E8%82%A1%E4%BB%BD%E6%9C%89%E9%99%90%E5%85%AC%E5%8F%B8%26num%3D100&exp=",
		"https://patents.google.com/xhr/query?url=inventor%3DJIN%2BLIU%26assignee%3DAlcatel%2BLucent%26num%3D100&exp=",
		"https://patents.google.com/xhr/query?url=inventor%3D%E5%88%98%E7%91%BE%26assignee%3D%E5%8D%8E%E4%B8%BA%E6%8A%80%E6%9C%AF%E6%9C%89%E9%99%90%E5%85%AC%E5%8F%B8%26num%3D100&exp=",
		"https://patents.google.com/xhr/query?url=inventor%3Dliu%2Bjin%26assignee%3DHuawei%2BTechnologies%2BCo.%252cLtd.%26num%3D100&exp=",
	}
	total := 0
	tmp := 0
	for _, url := range urls {
		tmp = getPatentTree(url, total+1)
		total = total + tmp
	}
}
