package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const captiveServerUrl = "http://www.google.cn/generate_204"
const serviceString = "在这里填写你的互联网接入商，需要对接入商的中文进行两次UrlEncode。如果无需选择互联网接入商，则此处留空。"

func printHelp() {
	fmt.Println("Usage: ./ruijie username password")
	fmt.Println("Example: ./ruijie 123456 123456")
}

func getCaptiveServerResponseStatusCodeAndBody() (int, string, error) {
	response, err := http.Get(captiveServerUrl)
	if err != nil {
		return -1, "", errors.New("can not send get request to captive server")
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return -1, "", errors.New("can not read captive server response body")
	}
	return response.StatusCode, string(body), nil
}

func getLoginUrlFromHtmlCode(htmlCode string) (string, string) {
	loginPageUrl := strings.Split(htmlCode, "'")[1]
	loginUrl := strings.Replace(strings.Split(loginPageUrl, "?")[0], "index.jsp", "InterFace.do?method=login", -1)
	queryString := strings.Replace(strings.Replace(strings.Split(loginPageUrl, "?")[1], "&", "%2526", -1), "=", "%253D", -1)
	return loginUrl, queryString
}

func login(loginUrl, username, password, serviceString, queryString string) (string, error) {
	client := &http.Client{}
	loginPostData := fmt.Sprintf("userId=%v&password=%v&service=%v&queryString=%v&operatorPwd=&operatorUserId=&validcode=&passwordEncrypt=false", username, password, serviceString, queryString)
	request, err := http.NewRequest(http.MethodPost, loginUrl, strings.NewReader(loginPostData))
	if err != nil {
		return "", errors.New("can not create login request")
	}
	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.132 Safari/537.36")
	response, err := client.Do(request)
	if err != nil {
		return "", errors.New("can not send login request")
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", errors.New("can not read login response body")
	}
	return string(body), nil
}

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}
	// Check network status
	captiveServerStatusCode, captiveServerResponseBody, err := getCaptiveServerResponseStatusCodeAndBody()
	if err != nil {
		log.Println(err.Error())
		return
	}
	if captiveServerStatusCode == 204 {
		// Exit when user is already online
		log.Println("You are already online!")
		return
	}
	// Start ruijie login
	loginUrl, queryString := getLoginUrlFromHtmlCode(captiveServerResponseBody)
	username := os.Args[1]
	password := os.Args[2]
	loginResult, err := login(loginUrl, username, password, serviceString, queryString)
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println(loginResult)
}
