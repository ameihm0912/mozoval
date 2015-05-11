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

func main() {
	var useRegexp bool

	flag.BoolVar(&useRegexp, "r", false, "argument is regular expression")
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "specify substring to match\n")
		os.Exit(1)
	}

	ret := oval.PackageQuery([]string{args[0]}, useRegexp)

	for _, x := range ret {
		fmt.Printf("%v %v %v\n", x.Name, x.Version, x.PkgType)
	}
}
