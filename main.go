package main

// TODO: Proxy support
// TODO: self-signed support
// TODO: auth support (basic/digest/ntlm?)
// TODO: cli arg: useragent
// TODO: cli arg: header
// TODO: cli arg: top x (only display top x words)
// TODO: exchange sync map with regular map + mutex
// TODO: sort wordlist by count
// TODO: on/offsite
// TODO: process response in extra thread via channel
// TODO: Cewl also extracts text from alt and title attributes (add this!)
// TODO: words from metadata
// TODO: E-Mails
// TODO: error handling
// TODO: support post login
// TODO: cobra cli
// TODO: Catch ctrl+c and write current map to wordlist

import (
	"github.com/kevin-ott/gocewl/cmd"
)

func main() {
	cmd.Start()

}
