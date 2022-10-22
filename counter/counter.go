package counter

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

type HttpGetter interface {
	Get(url string) (resp *http.Response, err error)
}

type counter struct {
	mx     sync.Mutex
	totals int

	client   HttpGetter
	maxProcs int
	timeout  time.Duration
}

var htmlRegex = regexp.MustCompile(`(?s)<script.*?>.*?<\/script>|<.*?>`)

func New(client HttpGetter, maxProcs int, timeout time.Duration) *counter {
	return &counter{
		client:   client,
		maxProcs: maxProcs,
		timeout:  timeout,
	}
}

func (c *counter) Count(word string, urls []string) int {
	limiter := make(chan int, c.maxProcs)

	var wg = new(sync.WaitGroup)
	wg.Add(len(urls))

	i := 0
	for {
		if i == len(urls) {
			close(limiter)
			break
		}
		url := urls[i]

		select {
		case limiter <- i:
			go func() {
				defer wg.Done()
				c.countForUrl(limiter, word, url)
			}()
			i++
		default:
			log.Println("limit")
			time.Sleep(c.timeout)
		}
	}

	wg.Wait()

	fmt.Printf("Total: %d", c.totals)
	return c.totals
}

// Updates total count with a number of words found in the given URL.
func (c *counter) countForUrl(done <-chan int, word, url string) {
	defer func() {
		<-done
	}()

	resp, err := c.client.Get(url)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	count := strings.Count(stripTags(string(body)), word)
	fmt.Printf("Count for %s: %d\n", url, count)

	c.mx.Lock()
	c.totals += count
	c.mx.Unlock()
}

func stripTags(s string) string {
	return htmlRegex.ReplaceAllString(s, " ")
}
