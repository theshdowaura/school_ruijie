package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const captiveServerUrl = "http://www.google.cn/generate_204"

func printHelp() {
	fmt.Println("Usage: ./ruijie service_name username password")
	fmt.Println("Service parameter can be \"chinamobile\", \"chinanet\" and \"chinaunicom\". If service parameter do not set as these value, it will use campus network as default internet service provider.")
	fmt.Println("Example: ./ruijie chinanet 201620000000 123456")
	fmt.Println("if you want to logout, use: ./ruijie logout")
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

func getServiceStringByServiceName(serviceName string) string {
	serviceStringMap := make(map[string]string)
	serviceStringMap["chinamobile"] = "%25E7%25A7%25BB%25E5%258A%25A8%25E5%25AE%25BD%25E5%25B8%25A6%25E6%258E%25A5%25E5%2585%25A5"
	serviceStringMap["chinanet"] = "%25E7%2594%25B5%25E4%25BF%25A1%25E5%25AE%25BD%25E5%25B8%25A6%25E6%258E%25A5%25E5%2585%25A5"
	serviceStringMap["chinaunicom"] = "%25E8%2581%2594%25E9%2580%259A%25E5%25AE%25BD%25E5%25B8%25A6%25E6%258E%25A5%25E5%2585%25A5"
	serviceStringMap["campus"] = "%25E6%2595%2599%25E8%2582%25B2%25E7%25BD%2591%25E6%258E%25A5%25E5%2585%25A5"
	if value, ok := serviceStringMap[serviceName]; ok {
		return value
	} else {
		return serviceStringMap["campus"]
	}
}

func getLoginUrlAndQueryStringFromHtmlCode(htmlCode string) (string, string) {
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

func getUserIndex() (string, error) {
	response, err := http.Get("http://10.8.2.2/eportal/redirectortosuccess.jsp")
	if err != nil {
		return "", errors.New("can not send get user index request")
	}
	defer response.Body.Close()
	userIndex := strings.Split(response.Request.URL.String(), "userIndex=")[1]
	return userIndex, nil
}

func logout() (string, error) {
	userIndex, err := getUserIndex()
	if err != nil {
		return "", err
	}
	response, err := http.PostForm("http://10.8.2.2/eportal/InterFace.do?method=logout", url.Values{
		"userIndex": {userIndex},
	})
	if err != nil {
		return "", errors.New("can not send logout request")
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", errors.New("can not read logout request body")
	}
	return string(body), nil
}

func main() {
	if len(os.Args) < 3 {
		if len(os.Args) == 2 && os.Args[1] == "logout" {
			logoutResult, err := logout()
			if err != nil {
				log.Println(err.Error())
			}
			log.Println(logoutResult)
		} else {
			// Show help when len(os.Args) < 3 and no logout
			printHelp()
		}
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
	loginUrl, queryString := getLoginUrlAndQueryStringFromHtmlCode(captiveServerResponseBody)
	serviceName := os.Args[1]
	username := os.Args[2]
	password := os.Args[3]
	serviceString := getServiceStringByServiceName(serviceName)
	loginResult, err := login(loginUrl, username, password, serviceString, queryString)
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println(loginResult)
}
