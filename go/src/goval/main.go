package main

import (
	"fmt"
	"flag"
	"os"
	"oval"
)

const (
	OPMODE_LIST = 1
)

var cfg Config

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
	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(2)
	}
	flag.Parse()

	if (cfg.flag_list == "path") {
		flag.Usage()
		os.Exit(2)
	} else {
		opmode = OPMODE_LIST
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
	default:
		flag.Usage()
	}
}
