package utilsEx

import (
	"io/ioutil"
	"net/http"
	"strings"
)

//var urls string

//method ： GET,POST,PUT,DELETE
//url : http url
//msg : json string
//headers: need  add header
func HttpSend(method string, url string, msg []byte, headers map[string]string) (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest(method, url, strings.NewReader(string(msg)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	for key, header := range headers {
		req.Header.Set(key, header)
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	// defer func() {
	// 	if resp != nil {
	// 		resp.Body.Close()
	// 	}
	// }()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

//设置为simple请求
func HttpSendfortext(method string, url string, msg []byte, headers map[string]string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	// defer func() {
	// 	if resp != nil {
	// 		resp.Body.Close()
	// 	}
	// }()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
