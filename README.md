# goCeWL version 0.2

Clone of digininja's [CeWL](https://github.com/digininja/CeWL) written in Golang.

- Crawl websites concurrently and extract words into a wordlist
- Should be faster as the original CeWL, as requests and parsing are performed concurrently. 
- static binary available, so no dependencies required
- lower memory fooprint

**Note**: This repo is experimental. Cosider it pre-alpha. The api / cli can change at any time.


## Installation

Note: Currently there are no tagged releases or pre-compiled binaries. This will change in the future. 

To compile and and install goCeWL, Go needs to be installed on your system. If that's not yet the case, please follow the installation instructions [here](https://golang.org/doc/install).

If you have Go installed, run `go get github.com/kevin-ott/gocewl`. This will download all dependencies and install the binary to `$GOPATH/bin`. 

## Usage
Run `gocewl --help` to display the commandline options.

```
gocewl is a commandline tool to generate custom wordlists by crawling webpages. It is based on CewL by digininja.

Usage:
  gocewl URL [flags]

Flags:
  -A, --allow stringArray   Domains in scope for the crawler. Provide as comma sperated list.
  -d, --depth int           Maximum depth for crawling (default 2)
  -h, --help                help for gocewl
  -k, --insecure            Ignore self-signed certificates
      --max-word int        Maximum word length (default 15)
  -c, --min-count int       Minimum number of times that the word was found (default 1)
      --min-word int        Mininum word length (default 3)
  -O, --offsite             Allow the crawler to visit offsite domains
  -p, --proxy string        Proxy to use: http[s]://[user:pass@]proxy.example.com[:8080]
  -q, --quiet               No output, except for words
  -t, --threads int         Amount of threads for crawling (default 10)
  -u, --url string          URL to start crawling
      --user-agent string   Custom user agent (default "gocewl/0.1")
      --version             version for gocewl
  -w, --write string        filename to write the wordlist to. If no file is provided, print to stdout (default "wordlist.txt")
```


## Examples
Crawl https://en.wikipedia.org with default parameters.

```gocewl <url>```

## Todos

### Parity with CeWL
- [x] Set minimum word length (defaults to 5)
- [x] Set crawling depth (defaults to 2)
- [x] Allow offsite crawling
- [x] Proxy support
- [ ] HTTP Basic / NTLM Auth support
- [ ] Include E-Mails
- [ ] Include metadata
- [ ] Headers 
- [x] User-agent

### Planned features
- [ ] Cookie support
- [ ] Sort wordlist by wordcount
- [ ] --top-words cli switch to only print top X words (by count)

### Other 
- [ ] Performance optimizations
- [ ] Improved error handling
- [ ] Improved cli

## Changelog

### 0.2
- Performance improvements
- Changed sync.Map to regular map with a mutex
- Fixed a race consdition when counting requests and error
- Fixed display of statistics

### 0.1
- Initial release to github