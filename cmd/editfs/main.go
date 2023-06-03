package main

import (
	"flag"
	"log"
	"os"

	"github.com/altid/editfs"
)

var (
	srv   = flag.String("s", "edit", "name of service")
	addr  = flag.String("a", "localhost", "listening address")
	mdns  = flag.Bool("m", false, "enable mDNS broadcast of service")
	port  = flag.Int("p", 12345, "default port to listen on")
	debug = flag.Bool("d", false, "enable debug printing")
	setup = flag.Bool("conf", false, "run configuration setup")
)

func main() {
	flag.Parse()
	if flag.Lookup("h") != nil {
		flag.Usage()
		os.Exit(1)
	}

	if *setup {
		if e := editfs.CreateConfig(*srv, *debug); e != nil {
			log.Fatal(e)
		}
		os.Exit(1)
	}

	edit, err := editfs.Register(*addr, *port, *srv, *debug)
	if err != nil {
		log.Fatal(err)
	}
	defer edit.Cleanup()
	if *mdns {
		if e := edit.Broadcast(); e != nil {
			log.Fatal(e)
		}
	}

	if e := edit.Run(); e != nil {
		log.Fatal(e)
	}

	os.Exit(0)
}
