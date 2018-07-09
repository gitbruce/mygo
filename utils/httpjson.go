package utils

func HttpJson(urlStr string, proxy bool) string {
	hc := new(HttpConnection)
	hc.Init(proxy)
	body, _ := hc.GetHttpBody(urlStr)
	return body
}

