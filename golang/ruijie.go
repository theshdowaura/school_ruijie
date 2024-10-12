package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-ping/ping"
)

const (
	captiveServerUrl = "http://www.google.cn/generate_204"
	pingHost         = "180.101.50.188" // 设置 Ping 的目标服务器
)

func getCaptiveServerResponseStatusCodeAndBody() (int, string, error) {
	response, err := http.Get(captiveServerUrl)
	if err != nil {
		return -1, "", errors.New("can not send get request to captive server")
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
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

func login(loginUrl, username, password, queryString, servicespasswd string) (string, error) {
	client := &http.Client{}
	loginPostData := fmt.Sprintf("userId=%v&password=%v&service=&queryString=%v&operatorPwd=%v&operatorUserId=&validcode=&passwordEncrypt=false", username, password, queryString, servicespasswd)
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
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", errors.New("can not read login response body")
	}
	return string(body), nil
}

func startPingCheck(username, password, servicespasswd string) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 使用 ping 检测网络连接
			pinger, err := ping.NewPinger(pingHost)
			if err != nil {
				log.Println("Failed to create pinger:", err)
				continue
			}
			pinger.Count = 3
			pinger.Timeout = time.Second * 10
			err = pinger.Run()
			if err != nil || pinger.Statistics().PacketLoss > 0 {
				log.Println("Ping failed, packet loss detected. Re-authenticating...")
				// 检查网络状态并重新认证
				err := reAuthenticate(username, password, servicespasswd)
				if err != nil {
					log.Println("Re-authentication failed:", err)
				} else {
					log.Println("Re-authentication successful!")
				}
			} else {
				log.Println("Network is stable.")
			}
		}
	}
}

func reAuthenticate(username, password, servicespasswd string) error {
	// Check network status
	captiveServerStatusCode, captiveServerResponseBody, err := getCaptiveServerResponseStatusCodeAndBody()
	if err != nil {
		return err
	}
	if captiveServerStatusCode == 204 {
		// Exit when user is already online
		log.Println("You are already online!")
		return nil
	}
	// Start ruijie login
	loginUrl, queryString := getLoginUrlFromHtmlCode(captiveServerResponseBody)
	loginResult, err := login(loginUrl, username, password, queryString, servicespasswd)
	if err != nil {
		return err
	}
	log.Println(loginResult)
	return nil
}

func main() {
	u := flag.String("u", "", "school_id")
	p := flag.String("p", "", "school_id passwd")
	c := flag.String("c", "", "school_id Carrier password")
	e := flag.String("e", "", "Persistent Login status")
	flag.Parse()

	username := *u
	password := *p
	servicespasswd := *c
	environment := *e
	// 初次认证
	err := reAuthenticate(username, password, servicespasswd)
	if err != nil {
		log.Println("Initial authentication failed:", err)
		return
	}

	// 开始 Ping 检测
	if environment == "on" {
		startPingCheck(username, password, servicespasswd)
	}

}
