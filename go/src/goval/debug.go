package main

import (
	"fmt"
	"os"
)

var debug_flag bool

func set_debug(f bool) {
	debug_flag = f
}

func debug_prt(s string, args ...interface{}) {
	if !debug_flag {
		return
	}
	fmt.Fprintf(os.Stderr, s, args...)
}
