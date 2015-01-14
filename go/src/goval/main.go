package main

import (
	"fmt"
	"flag"
	"os"
	"oval"
)

const (
	OPMODE_LIST = 1
	OPMODE_RUN = 2
)

var cfg config

func run_mode() {
	od, ret := oval.Parse(cfg.flag_run)
	if ret != nil {
		fmt.Fprintf(os.Stderr, "%v\n", ret)
	}
	oval.Execute(od)
}

func list_mode() {
	od, ret := oval.Parse(cfg.flag_list)
	if ret != nil {
		fmt.Fprintf(os.Stderr, "%v\n", ret)
	}

	for _, v := range od.Definitions.Definitions {
		fmt.Printf("%s %s\n", v.ID, v.Metadata.Title)
	}
}

func main() {
	var opmode int = 0

	cfg = default_config()
	flag.BoolVar(&cfg.flag_debug, "d", false, "enable debugging")
	flag.StringVar(&cfg.flag_list, "l", "path", "list checks")
	flag.StringVar(&cfg.flag_run, "r", "path", "run checks")
	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(2)
	}
	flag.Parse()

	var validmode bool = false
	if (cfg.flag_list != "path") {
		opmode = OPMODE_LIST
		validmode = true
	} else if (cfg.flag_run != "path") {
		opmode = OPMODE_RUN
		validmode = true
	}
	if !validmode {
		flag.Usage()
		os.Exit(2)
	}

	oval.Init()

	if (cfg.flag_debug) {
		set_debug(true)
		// If we enable debugging on the command line, also turn it on
		// in the OVAL library
		oval.Set_debug(true)
		debug_prt("Debugging enabled\n")
	}

	switch opmode {
	case OPMODE_LIST:
		debug_prt("Entering list mode\n")
		list_mode()
	case OPMODE_RUN:
		debug_prt("Entering run mode\n")
		run_mode()
	default:
		flag.Usage()
	}
}
