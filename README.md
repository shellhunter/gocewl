# goCeWL

Clone of digininja's CeWL written in Golang.

**Note**: This repo is experimental. Cosider it pre-alpha. The api / cli can change at any time.

## Why?
* I hate ruby
* I like go
* I like cewl

**(in all seriousness)**
* Should be faster, as requests and parsing are performed concurrently. 
* static binary available, so no dependencies required
* lower memory fooprint

## Features
* Crawl websites concurrently and extract words into a wordlist

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