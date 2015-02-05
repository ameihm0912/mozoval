package main

import (
	"fmt"
	"os"
)

var debugFlag bool

func setDebug(f bool) {
	debugFlag = f
}

func debugPrint(s string, args ...interface{}) {
	if !debugFlag {
		return
	}
	fmt.Fprintf(os.Stdout, s, args...)
}
