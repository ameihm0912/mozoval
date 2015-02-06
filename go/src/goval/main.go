// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Contributor:
// - Aaron Meihm ameihm@mozilla.com
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

func runMode() {
	od, ret := oval.Parse(cfg.flagRun)
	if ret != nil {
		fmt.Fprintf(os.Stderr, "%v\n", ret)
		os.Exit(1)
	}

	results := oval.Execute(od)
	for _, v := range results {
		notice := "--"
		if v.StatusString() != "false" {
			notice = "!!"
		}
		fmt.Fprintf(os.Stdout, "%v %v %v %v\n", notice, v.ID, v.StatusString(), v.Title)
	}
}

func listMode() {
	od, ret := oval.Parse(cfg.flagList)
	if ret != nil {
		fmt.Fprintf(os.Stderr, "%v\n", ret)
		os.Exit(1)
	}

	for _, v := range od.Definitions.Definitions {
		fmt.Printf("%v %v\n", v.ID, v.Metadata.Title)
	}
}

func main() {
	var opmode int = 0

	cfg = defaultConfig()
	flag.BoolVar(&cfg.flagDebug, "d", false, "enable debugging")
	flag.StringVar(&cfg.flagList, "l", "path", "list checks")
	flag.StringVar(&cfg.flagRun, "r", "path", "run checks")
	flag.IntVar(&cfg.maxChecks, "n", 10, "concurrent checks")
	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(2)
	}
	flag.Parse()

	var validmode bool = false
	if cfg.flagList != "path" {
		opmode = OPMODE_LIST
		validmode = true
	} else if cfg.flagRun != "path" {
		opmode = OPMODE_RUN
		validmode = true
	}
	if !validmode {
		flag.Usage()
		os.Exit(2)
	}

	oval.Init()

	if cfg.flagDebug {
		setDebug(true)
		// If we enable debugging on the command line we also enable
		// it in the OVAL library.
		oval.SetDebug(true)
		debugPrint("debugging enabled\n")
	}
	oval.SetMaxChecks(cfg.maxChecks)

	switch opmode {
	case OPMODE_LIST:
		debugPrint("entering list mode\n")
		listMode()
	case OPMODE_RUN:
		debugPrint("entering run mode\n")
		runMode()
	default:
		flag.Usage()
	}
}
