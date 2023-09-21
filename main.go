package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Website struct {
	Name     string
	Url      string
	Snapshot string
}

var interval = 0 // in seconds
var websites = []Website{}

func main() {
	configuration := readConfiguration()
	initializeWebsites(configuration)
	goToSleep()

	for {
		for i := 0; i < len(websites); i++ {
			websites[i] = checkWebsite(websites[i])
		}
		goToSleep()
	}
}

func initializeWebsites(configuration Configuration) {
	websites = configuration.Websites
	interval = configuration.Interval

	for i := 0; i < len(websites); i++ {
		websites[i] = createInitialSnapshot(websites[i])
	}
}

func createInitialSnapshot(website Website) Website {
	content, error := getWebsiteAsString(website)
	if error == nil {
		website.Snapshot = content
	}
	log.Println("Created initial snapshot for " + website.Name)
	return website
}

func getWebsiteAsString(website Website) (string, error) {
	resp, err := http.Get(website.Url)
	if err != nil {
		log.Println(err)
		log.Println("An error occurred! The website could not be reached!")
		return "Error", errors.New("The website could not be reached")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	content := string(body[:])
	content = sanitizeHtml(content)
	return content, nil
}

func checkWebsite(website Website) Website {
	content, error := getWebsiteAsString(website)
	if error != nil {
		return website
	}
	if website.Snapshot != content {
		website.Snapshot = content
		printContentChangeMsg(website)
		playSound()
	} else {
		log.Println("No changes for " + website.Name)
	}

	return website
}

func printContentChangeMsg(website Website) {
	log.Println("========= " + website.Name + " =========")
	log.Println("Content changed: " + website.Url)
	log.Println("====================" + strings.Repeat("=", len(website.Name)))
}

func goToSleep() {
	log.Printf("Going to sleep for %s seconds", strconv.FormatInt(int64(interval), 10))
	log.Println()
	time.Sleep(time.Duration(interval) * time.Second)
}
