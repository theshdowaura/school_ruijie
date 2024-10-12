package main

import (
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
	captiveServerURL = "http://www.google.cn/generate_204"
	pingHost         = "180.101.50.188" // Ping target server
)

func getCaptiveServerResponseStatusCodeAndBody() (int, string, error) {
	response, err := http.Get(captiveServerURL)
	if err != nil {
		return -1, "", fmt.Errorf("failed to send GET request to captive server: %w", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return -1, "", fmt.Errorf("failed to read captive server response body: %w", err)
	}
	return response.StatusCode, string(body), nil
}

func getLoginURLFromHTMLCode(htmlCode string) (string, string) {
	loginPageURL := strings.Split(htmlCode, "'")[1]
	loginURL := strings.Replace(strings.Split(loginPageURL, "?")[0], "index.jsp", "InterFace.do?method=login", -1)
	queryString := strings.Replace(strings.Replace(strings.Split(loginPageURL, "?")[1], "&", "%2526", -1), "=", "%253D", -1)
	return loginURL, queryString
}

func login(loginURL, username, password, queryString, servicesPassword string) (string, error) {
	client := &http.Client{}
	loginPostData := fmt.Sprintf(
		"userId=%v&password=%v&service=&queryString=%v&operatorPwd=%v&operatorUserId=&validcode=&passwordEncrypt=false",
		username, password, queryString, servicesPassword,
	)

	request, err := http.NewRequest(http.MethodPost, loginURL, strings.NewReader(loginPostData))
	if err != nil {
		return "", fmt.Errorf("failed to create login request: %w", err)
	}

	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.132 Safari/537.36")

	response, err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("failed to send login request: %w", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read login response body: %w", err)
	}
	return string(body), nil
}

func startPingCheck(username, password, servicesPassword string) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		pinger, err := ping.NewPinger(pingHost)
		if err != nil {
			log.Printf("failed to create pinger: %v", err)
			continue
		}

		pinger.Count = 3
		pinger.Timeout = 10 * time.Second
		err = pinger.Run()
		if err != nil || pinger.Statistics().PacketLoss > 0 {
			log.Println("Ping failed, packet loss detected. Re-authenticating...")
			if err := reAuthenticate(username, password, servicesPassword); err != nil {
				log.Printf("re-authentication failed: %v", err)
			} else {
				log.Println("Re-authentication successful!")
			}
		} else {
			log.Println("Network is stable.")
		}
	}
}

func reAuthenticate(username, password, servicesPassword string) error {
	statusCode, responseBody, err := getCaptiveServerResponseStatusCodeAndBody()
	if err != nil {
		return err
	}

	if statusCode == http.StatusNoContent {
		log.Println("You are already online!")
		return nil
	}

	loginURL, queryString := getLoginURLFromHTMLCode(responseBody)
	loginResult, err := login(loginURL, username, password, queryString, servicesPassword)
	if err != nil {
		return err
	}

	log.Println(loginResult)
	return nil
}

func main() {
	username := flag.String("u", "", "school_id")
	password := flag.String("p", "", "school_id password")
	servicesPassword := flag.String("c", "", "carrier password")
	environment := flag.String("e", "", "persistent login status")
	flag.Parse()

	if *username == "" || *password == "" || *servicesPassword == "" {
		log.Fatal("username, password, and carrier password must be provided")
	}

	// Initial authentication
	if err := reAuthenticate(*username, *password, *servicesPassword); err != nil {
		log.Printf("initial authentication failed: %v", err)
		return
	}

	// Start Ping check if environment is "on"
	if *environment == "on" {
		startPingCheck(*username, *password, *servicesPassword)
	}
}
