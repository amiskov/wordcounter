package main

import (
	"net/http"
	"time"
	"wordcounter/counter"
)

func main() {
	urls := []string{
		"https://go.dev",
		"https://golangcourse.ru",
		"https://go.dev", // duplicates count
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	c := counter.New(client, 5, 1*time.Second)
	c.Count("Go", urls)
}
