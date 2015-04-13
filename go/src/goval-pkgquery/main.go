// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Contributor:
// - Aaron Meihm ameihm@mozilla.com
package main

import (
	"fmt"
	"os"
	"oval"
)

func main() {
	oval.Init()

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "specify substring to match\n")
		os.Exit(1)
	}

	args := make([]string, 0)
	args = append(args, os.Args[1])
	ret := oval.PackageQuery(args)

	for _, x := range ret {
		fmt.Printf("%v %v\n", x.Name, x.Version)
	}
}
