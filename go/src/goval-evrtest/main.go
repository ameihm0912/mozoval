package main

import (
	"bufio"
	"fmt"
	"os"
	"oval"
	"strings"
)

func main() {
	oval.Init()

	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, "Specify input test data as argument\n")
		os.Exit(1)
	}
	fmt.Println("starting evr comparison tests")

	oval.Set_debug(true)

	fd, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		buf := strings.TrimSpace(scanner.Text())
		if len(buf) == 0 {
			continue
		}
		fmt.Printf("%v\n", buf)

		var opmode int
		s0 := strings.Fields(buf)
		switch s0[0] {
		case "=":
			opmode = oval.EVROP_EQUALS
		case "<":
			opmode = oval.EVROP_LESS_THAN
		default:
			fmt.Fprintf(os.Stderr, "Unknown operation %v\n", s0[0])
			os.Exit(1)
		}
		result := oval.Test_evr_compare(opmode, s0[1], s0[2])
		if !result {
			fmt.Println("FAIL")
			os.Exit(2)
		}
		fmt.Println("PASS")
	}
	fd.Close()

	fmt.Println("end evr comparison tests")
}
