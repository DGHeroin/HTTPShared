package main

import (
	"log"
	"flag"
	"github.com/DGHeroin/HTTPShared"
	"os"
)

func parseArgs() (string, string) {
	Addr  := *flag.String("a", ":9999", "Bind Address")
	Token := *flag.String("t", "", "Auth Token")
	Help  := *flag.Bool  ("h", false, "Show Usage")
	flag.Parse()

	if Help {
		flag.Usage()
		os.Exit(0)
	}

	return Addr, Token
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	Addr, Token := parseArgs()

	web := HTTPShared.NewWebService(Addr, Token)
	web.WaitExit()
}
