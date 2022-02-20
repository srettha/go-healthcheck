package checker

import (
	"challenge/go-healthcheck/client"
	"sync"
	"time"
)

type Result struct {
	TotalWebsites int   `json:"total_websites"`
	Success       int   `json:"success"`
	Failure       int   `json:"failure"`
	TotalTime     int64 `json:"total_time"`
}

// Ping url with http client
func PingUrl(client client.HttpClient, url string) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}

// Format ping result
func FormatPingResult(totalWebsites int, success int, failure int, elapse time.Duration) Result {
	totalTime := int64(elapse / time.Millisecond)

	return Result{
		TotalWebsites: totalWebsites,
		Success:       success,
		Failure:       failure,
		TotalTime:     totalTime,
	}
}

// Ping urls
func Ping(client client.HttpClient, urls []string) Result {
	start := time.Now()

	var success, failure int
	var wg sync.WaitGroup

	for _, u := range urls {
		wg.Add(1)

		go func(url string) {
			defer wg.Done()

			err := PingUrl(client, url)
			if err == nil {
				success += 1
			} else {
				failure += 1
			}
		}(u)
	}

	wg.Wait()

	elapse := time.Since(start)

	return FormatPingResult(len(urls), success, failure, elapse)
}
