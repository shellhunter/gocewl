package gocewl

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"

	query "github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

var (
	extRegex      = regexp.MustCompile(`(\.zip$|\.gz$|\.zip$|\.bz2$|\.png$|\.gif$|\.jpg$|^#)`)
	newlineRegex  = regexp.MustCompile(`\n+`)       //one or more newlines
	blankRegex    = regexp.MustCompile(`\s{2,}`)    //two or more spaces
	notAlphaRegex = regexp.MustCompile("[^a-zA-Z]") // all non-alpha characters
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
func prepareSignalHandler() {}
func sortMap()              {}
func printBanner()          {}
func printConfig()          {}

type Stats struct {
	RequestCount  int
	ResponseCount int
	ErrorCount    int
	TotalWords    int
	TotalTime     int
}

func (s *Stats) String() string {
	return fmt.Sprintf(`
		Total requests: %d
		Total responses: %d
		Total Errors: %d
		Total Words: %d
		Total Time: %d`, s.RequestCount, s.ResponseCount, s.ErrorCount, s.TotalWords, s.TotalTime)
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
	var stats Stats
	wordsWithCount := sync.Map{}

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
	}

	// Callbacks

	crawler.OnError(func(r *colly.Response, err error) {
		stats.ErrorCount++
		fmt.Printf("[!] Request to URL %s failed. Reason: %s\n", r.Request.URL.String(), err.Error())
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

	crawler.OnRequest(func(_ *colly.Request) {
		stats.RequestCount++
	})

	crawler.OnResponse(func(r *colly.Response) {
		stats.ResponseCount++

		words := extractFromText(r.Body)
		stats.TotalWords += len(words)
		if !config.Quiet {
			fmt.Printf("[%d] --> %s (%d)\n", r.StatusCode, r.Request.URL, len(words))
		}
		for _, word := range words {
			word = notAlphaRegex.ReplaceAllString(word, "")
			if len(word) < config.MininumWordLength || len(word) > config.MaximumWordLength {
				continue
			}
			c, ok := wordsWithCount.Load(word)
			if ok {
				count := c.(int)
				count++
				wordsWithCount.Store(word, count)
			} else {
				wordsWithCount.Store(word, 1)
			}
		}
	})
	crawler.Visit(seedURL.String())
	crawler.Wait()

	fh, err := os.Create(config.OutputFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()

	wordsWithCount.Range(func(key, value interface{}) bool {
		word := key.(string)
		count := value.(int)
		if count >= config.MinimumWordCount {
			fmt.Printf("%d\t%s\n", count, word)
			_, err := fh.WriteString(word + "\n")
			if err != nil {
				log.Fatal(err)
			}
		}
		return true
	})
	fmt.Println(stats)
}
