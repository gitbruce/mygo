package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type HttpConnection struct {
	// UserAgent is the User-Agent string used by HTTP requests
	UserAgent string
	Proxy bool
}

const timeout time.Duration = 400

func (h *HttpConnection) ua(ua string) {
	h.UserAgent = ua
}

func (h *HttpConnection) Init(useProxy bool) {
	h.ua("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36")
	h.Proxy = useProxy
}

func (h *HttpConnection) GetHttpClient() *http.Client {
	var client http.Client
	var proxy string
	if h.Proxy {
		proxy = "socks5://127.0.0.1:1086"
	}
	if proxy != "" {
		proxyURL, err := url.Parse(proxy)
		if err != nil {
			panic("Error parsing proxy URL")
		}
		transport := http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
		client = http.Client{
			Transport: &transport,
		}
	} else {
		client = http.Client{}
	}
	return &client

}

func (h *HttpConnection) GetHttpResponse(urlStr string) (*http.Response, error) {
	client := h.GetHttpClient()
	response, err := client.Get(urlStr)
	return response, err
}



func (h *HttpConnection) GetHttpBody(urlStr string) (string, error) {
	response, err := h.GetHttpResponse(urlStr)
	defer response.Body.Close()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		body, err := ioutil.ReadAll(response.Body)
		return string(body), err
	}
	return "", err
}

//func main() {
//	hc := new(HttpConnection)
//	hc.UserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36"
//	s := hc.get("http://www.google.com", "socks5://127.0.0.1:1086")
//	fmt.Println(s)
//}
