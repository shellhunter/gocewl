# goCeWL

Clone of digininja's [CeWL](https://github.com/digininja/CeWL) written in Golang.

**Note**: This repo is experimental. Cosider it pre-alpha. The api / cli can change at any time.

- Crawl websites concurrently and extract words into a wordlist
- Should be faster as the original CeWL, as requests and parsing are performed concurrently. 
- static binary available, so no dependencies required
- lower memory fooprint

## Features
- Crawl websites concurrently and extract words into a wordlist

## Installation

Note: Currently there are no tagged releases or pre-compiled binaries. This will change in the future. 

To compile and and install goCeWL, Go needs to be installed on your system. If that's not yet the case, please follow the installation instructions [here](https://golang.org/doc/install).

If you have Go installed, run `go get github.com/kevin-ott/gocewl`. This will download all dependencies and install the binary to `$GOPATH/bin`. 

## Usage
Run `gocewl --help` to display the commandline options.

`gocewl -u <url>` will crawl the specified URL with the default parameters. 

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
- [ ] JS parsing with headless chrome
- [ ] Cookie support
- [ ] Authenticated crawling
- [ ] Login sequences with yaml-files
- [ ] Sort wordlist by wordcount
- [ ] --top-words cli switch to only print top X words (by count)
### Other 
- [ ] Performance optimizations
- [ ] Improved error handling
- [ ] Improved cli

## Changelog

### 0.1
- Initial release to github