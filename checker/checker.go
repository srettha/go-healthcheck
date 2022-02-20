package checker

import (
	"challenge/go-healthcheck/client"
	"time"
)

type Result struct {
	TotalWebsites int   `json:"total_websites"`
	Success       int   `json:"success"`
	Failure       int   `json:"failure"`
	TotalTime     int64 `json:"total_time"`
}

type PingUrlResult struct {
	err error
}

// Ping url with http client
func PingUrl(client client.HttpClient, url string) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}

	resp.Body.Close()

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
	totalWebsites := len(urls)
	start := time.Now()

	var success, failure int
	resultsChan := make(chan *PingUrlResult)

	defer close(resultsChan)

	for _, url := range urls {
		go func(url string) {
			err := PingUrl(client, url)
			result := &PingUrlResult{err}
			resultsChan <- result
		}(url)
	}

	var results []PingUrlResult
	for result := range resultsChan {
		results = append(results, *result)
		if result.err == nil {
			success += 1
		} else {
			failure += 1
		}

		if len(results) == totalWebsites {
			break
		}
	}

	elapse := time.Since(start)

	return FormatPingResult(totalWebsites, success, failure, elapse)
}
