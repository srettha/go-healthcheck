package main

import (
	"challenge/go-healthcheck/checker"
	"challenge/go-healthcheck/client"
	"challenge/go-healthcheck/oauth"
	"challenge/go-healthcheck/reader"
	"challenge/go-healthcheck/reporter"
	"log"
	"os"
)

func main() {
	channelID := os.Getenv("LINE_CHANNEL_ID")
	channelSecret := os.Getenv("LINE_CHANNEL_SECRET")
	if channelID == "" || channelSecret == "" {
		log.Fatal("Channel ID and Channel secret are required")
	}

	baseRedirectURL := os.Getenv("BASE_REDIRECT_URL")
	if baseRedirectURL == "" {
		baseRedirectURL = "http://localhost:5555"
		log.Printf("Base redirect url has been set to %s", baseRedirectURL)
	}

	filePath := os.Args[1]
	if filePath == "" {
		log.Fatal("Input file is required")
	}

	urls, err := reader.OpenAndReadFile(filePath)
	if err != nil {
		log.Fatal("Unable to read input file")
	}

	log.Println("Perform website checking...")

	httpClient := client.GetHttpClient()
	pingResult := checker.Ping(httpClient, urls)

	socialClient := client.GetSocialClient(channelID, channelSecret)
	if err != nil {
		log.Fatal("Social SDK:", socialClient, " err:", err)
	}

	if err := oauth.LoginUser(socialClient, baseRedirectURL); err != nil {
		log.Fatal("Unable to login user:", err)
	}

	accessToken := oauth.AuthorizeUser(socialClient, baseRedirectURL)

	if err := reporter.Report(httpClient, accessToken, pingResult); err != nil {
		log.Println("Failed to send report to Healcheck Report system")
	}

	log.Println("Done!")
	log.Println("")

	log.Println("Checked website(s):", pingResult.TotalWebsites)
	log.Println("Successful website(s):", pingResult.Success)
	log.Println("Failure website(s):", pingResult.Failure)
	log.Println("Total times to finished checking website(s):", pingResult.TotalTime)
}
