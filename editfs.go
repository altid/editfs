package main

import (
	"flag"
	"log"
	"os"

	"github.com/altid/libs/config"
	"github.com/altid/libs/config/types"
	"github.com/altid/libs/fs"
)

var (
	mtpt    = flag.String("p", "/tmp/altid", "Path for filesystem")
	srv     = flag.String("s", "edit", "Name of service")
	cfgfile = flag.String("c", "", "Directory of configuration file")
	debug   = flag.Bool("d", false, "enable debug logging")
	setup   = flag.Bool("conf", false, "Run configuration setup")
)

func main() {
	flag.Parse()
	if flag.Lookup("h") != nil {
		flag.Usage()
		os.Exit(1)
	}

	conf := &struct {
		Listen types.ListenAddress `altid:"listen_address,no_prompt"`
		Logdir types.Logdir        `altid:"logdir,no_prompt"`
	}{"none", "none"}

	if *setup {
		if e := config.Create(conf, *srv, *cfgfile, *debug); e != nil {
			log.Fatal(e)
		}

		os.Exit(0)
	}

	if e := config.Marshal(conf, *srv, *cfgfile, *debug); e != nil {
		log.Fatal(e)
	}

	s := &server{}

	ctrl, err := fs.New(s, string(conf.Logdir), *mtpt, *srv, "feed", *debug)
	if err != nil {
		log.Fatal(err)
	}

	defer ctrl.Cleanup()
	ctrl.SetCommands(Commands...)
	ctrl.CreateBuffer("main", "feed")

	//go s.listen()

	if e := ctrl.Listen(); e != nil {
		log.Fatal(e)
	}
}
