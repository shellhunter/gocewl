package cmd

import (
	"fmt"
	"os"
	"strconv"

	cwl "github.com/kevin-ott/gocewl/crawler"
	"github.com/spf13/cobra"
)

var rootCMD = cobra.Command{
	Use:     "gocewl",
	Short:   "gocewl",
	Long:    `gocewl is a commandline tool to generate custom wordlists by crawling webpages. It is based on CewL by digininja.`,
	Version: "0.1",
	Run:     startCrawling,
}

var conf cwl.Config

func init() {

	rootCMD.Flags().BoolVarP(&conf.SkipSSL, "insecure", "k", false, "Ignore self-signed certificates")
	rootCMD.Flags().BoolP("quiet", "q", false, "No output, except for words")
	rootCMD.Flags().BoolVarP(&conf.Offsite, "offsite", "O", false, "Allow the crawler to visit offsite domains")

	//rootCMD.Flags().Int("top", 0, "Print n words with the highest count. If 0, all words are printed / written")
	rootCMD.Flags().IntVar(&conf.MininumWordLength, "min-word", 3, "Mininum word length")
	rootCMD.Flags().IntVar(&conf.MaximumWordLength, "max-word", 15, "Maximum word length")
	rootCMD.Flags().IntVarP(&conf.Depth, "depth", "d", 2, "Maximum depth for crawling")
	rootCMD.Flags().IntVarP(&conf.Threads, "threads", "t", 10, "Amount of threads for crawling")
	rootCMD.Flags().IntVarP(&conf.MinimumWordCount, "min-count", "c", 1, "Minimum number of times that the word was found")

	rootCMD.Flags().StringVar(&conf.UserAgent, "user-agent", "gocewl/0.1", "Custom user agent")
	rootCMD.Flags().StringVarP(&conf.URL, "url", "u", "", "URL to start crawling")
	rootCMD.Flags().StringVarP(&conf.OutputFilename, "write", "w", "wordlist.txt", "filename to write the wordlist to. If no file is provided, print to stdout")
	rootCMD.Flags().StringVarP(&conf.Proxy, "proxy", "p", "", "Proxy to use: http[s]://[user:pass@]proxy.example.com[:8080]")
	//rootCMD.Flags().String("auth-type", "", "")
	//rootCMD.Flags().String("auth-user", "", "")
	//rootCMD.Flags().String("auth-pass", "", "")

	//rootCMD.Flags().StringArray("headers", []string{}, "")
	rootCMD.Flags().StringArrayVarP(&conf.Domains, "allow", "A", []string{}, "Domains in scope for the crawler. Provide as comma sperated list.")

	rootCMD.MarkFlagRequired("url")
}

func startCrawling(cmd *cobra.Command, args []string) {
	// print config

	conf.URL = cmd.Flags().Lookup("url").Value.String()

	fmt.Println("url: " + conf.URL)
	// verify url scheme
	fmt.Println("depth: " + strconv.Itoa(conf.Depth))
	fmt.Println("minwordlength: " + strconv.Itoa(conf.MininumWordLength))
	fmt.Println("maxwordlength: " + strconv.Itoa(conf.MaximumWordLength))
	cwl.Crawl(&conf)
}

func Start() {
	if err := rootCMD.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
