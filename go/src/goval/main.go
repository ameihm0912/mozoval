package main

import (
	"flag"
	"fmt"
	"os"
	"oval"
)

const (
	_ = iota
	OPMODE_LIST
	OPMODE_RUN
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
	flag.IntVar(&cfg.max_checks, "n", 10, "concurrent checks")
	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(2)
	}
	flag.Parse()

	var validmode bool = false
	if cfg.flag_list != "path" {
		opmode = OPMODE_LIST
		validmode = true
	} else if cfg.flag_run != "path" {
		opmode = OPMODE_RUN
		validmode = true
	}
	if !validmode {
		flag.Usage()
		os.Exit(2)
	}

	oval.Init()

	if cfg.flag_debug {
		set_debug(true)
		// If we enable debugging on the command line, also turn it on
		// in the OVAL library
		oval.Set_debug(true)
		debug_prt("Debugging enabled\n")
	}
	oval.Set_max_checks(cfg.max_checks)

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
