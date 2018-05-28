package crawler

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"sync/atomic"
	"time"

	query "github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

var (
	extRegex      = regexp.MustCompile(`(\.zip$|\.gz$|\.zip$|\.bz2$|\.png$|\.gif$|\.jpg$|^#)`)
	newlineRegex  = regexp.MustCompile(`\n+`)       //one or more newlines
	blankRegex    = regexp.MustCompile(`\s{2,}`)    //two or more spaces
	notAlphaRegex = regexp.MustCompile("[^a-zA-Z]") // all non-alpha characters
	kill          = false
)

func extractFromAttributes() []string { return nil }
func extractFromText(body []byte) []string {
	dom, _ := query.NewDocumentFromReader(bytes.NewReader(body))
	dom.Find("script").Remove()
	dom.Find("style").Remove()
	text := dom.Text()
	text = newlineRegex.ReplaceAllString(text, " ")
	text = blankRegex.ReplaceAllString(text, " ")
	words := strings.Split(text, " ")
	return words
}
func prepareSignalHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	fmt.Println("Handling Keyboard interrupt.")
	kill = true
}
func printBanner()          {}
func printConfig()          {}
func writeResultsToFile()   {}
func writeResultsToStdout() {}
func resultsHandler()       {}

type Stats struct {
	RequestCount  uint64
	ResponseCount uint64
	ErrorCount    uint64
	TotalWords    uint64
	TotalTime     float64
}

func (s *Stats) String() string {
	return fmt.Sprintf(`
Total requests: %d
Total responses: %d
Total Errors: %d
Total Words: %d
Total Time: %.2fs`, s.RequestCount, s.ResponseCount, s.ErrorCount, s.TotalWords, s.TotalTime)
}

type Config struct {
	URL               string
	Depth             int
	SkipSSL           bool
	Threads           int
	Proxy             string
	UserAgent         string
	OutputPath        string
	OutputFilename    string
	MinimumWordCount  int
	MininumWordLength int
	MaximumWordLength int
	Quiet             bool
	Domains           []string
	Offsite           bool
	ShowCount         bool
}

func Crawl(config *Config) {
	go prepareSignalHandler()
	var stats Stats
	var wordsWithCount = NewWordMap()
	startTime := time.Now()

	seedURL, err := url.Parse(config.URL)
	if err != nil {
		log.Fatal(err)
	}

	crawler := colly.NewCollector(
		colly.MaxDepth(config.Depth),
		colly.Async(true),
		colly.UserAgent(config.UserAgent),
	)

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: config.SkipSSL},
	}

	crawler.WithTransport(transport)

	if config.Proxy != "" {
		proxyURL, err := url.Parse(config.Proxy)
		if err != nil {
			log.Fatal(err)
		}
		crawler.SetProxy(proxyURL.String())
	}

	// Threads / Delays

	crawler.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: config.Threads})
	if !config.Offsite {
		crawler.AllowedDomains = append(crawler.AllowedDomains, seedURL.Host)
	}

	if len(config.Domains) > 0 {
		crawler.AllowedDomains = append(crawler.AllowedDomains, seedURL.Host)
		crawler.AllowedDomains = append(crawler.AllowedDomains, config.Domains...)
	} else {
		crawler.AllowedDomains = append(crawler.AllowedDomains, seedURL.Host)
	}

	// Callbacks

	crawler.OnError(func(r *colly.Response, err error) {
		atomic.AddUint64(&stats.ErrorCount, 1)
		if !config.Quiet {
			fmt.Printf("[!] Request to URL %s failed. Reason: %s\n", r.Request.URL.String(), err.Error())
		}
	})

	crawler.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// check if the file is in the list of
		if bl := extRegex.FindString(link); bl != "" {
			return
		}
		absLink := e.Request.AbsoluteURL(link)
		e.Request.Visit(absLink)
	})

	crawler.OnRequest(func(r *colly.Request) {
		atomic.AddUint64(&stats.RequestCount, 1)
		if kill == true {
			r.Abort()
		}
	})

	crawler.OnResponse(func(r *colly.Response) {
		atomic.AddUint64(&stats.ResponseCount, 1)

		words := extractFromText(r.Body)
		atomic.AddUint64(&stats.TotalWords, uint64(len(words)))
		if !config.Quiet {
			fmt.Printf("[%d] --> %s (%d)\n", r.StatusCode, r.Request.URL, len(words))
		}
		for _, word := range words {
			word = notAlphaRegex.ReplaceAllString(word, "")
			if len(word) < config.MininumWordLength || len(word) > config.MaximumWordLength {
				continue
			}
			wordsWithCount.Add(word)
		}
	})
	if err := crawler.Visit(seedURL.String()); err != nil {
		log.Fatal(err)
	}
	crawler.Wait()

	fh, err := os.Create(config.OutputFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()

	for word, count := range wordsWithCount.internal {
		if count >= config.MinimumWordCount {
			//fmt.Printf("%s\n", word)
			_, err := fh.WriteString(word + "\n")
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	stats.TotalTime = time.Since(startTime).Seconds()
	fmt.Println(stats.String())
}
