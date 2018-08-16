package main

import (
	"HTTPShared"
	"log"
	"flag"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	Addr  := *flag.String("a", ":9999", "Bind Address")
	Token := *flag.String("t", "", "Auth Token")
	Help  := *flag.Bool  ("h", false, "Show Usage")
	flag.Parse()

	if Help {
		flag.Usage()
		return
	}

	web := HTTPShared.NewWebService(Addr, Token)
	web.WaitExit()
}
