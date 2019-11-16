package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)


//发送post请求
func HttpPost(url, postData string) (code int , str string) {

	resp, err := http.Post(url,
		"application/json",
		strings.NewReader(postData))
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		// handle error
		fmt.Println( "http 请求出现错误", err)
		return -1 ,""
	}

	fmt.Println(string(body))

	return resp.StatusCode, string(body)
}

//发送get请求
func HttpGet(url string) (code int , str string) {
	resp, err :=   http.Get(url)
	if err != nil {
		// handle error
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		return  -1, ""

	}
	fmt.Println(string(body))

	return resp.StatusCode, string(body)
}