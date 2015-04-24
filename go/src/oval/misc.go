// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Contributor:
// - Aaron Meihm ameihm@mozilla.com
package oval

import (
	"bufio"
	"io/ioutil"
	"os"
	"regexp"
)

func versionPtrnMatch(ver string, pattern string) bool {
	debugPrint("[versionPtrnMatch] %v ? %v\n", ver, pattern)
	// XXX Should handle errors here as the pattern can come from the
	// state as part of a definition.
	res, _ := regexp.MatchString(pattern, ver)
	return res
}

// Given a file, read the entire file and match against pattern. If
// we find a match, return it. If there are submatches that are part of
// the supplied pattern, we return the first submatch. A zero value
// string is returned if nothing is found.
func fileContentMatchAll(path string, pattern string) (ret string) {
	bytebuf, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return
	}

	subs := re.FindStringSubmatch(string(bytebuf))
	if len(subs) >= 2 {
		ret = subs[1]
	} else if len(subs) == 1 {
		ret = subs[0]
	}
	return
}

// Given a file, read the file line by line matching against pattern; if
// we find a match, return it. If there are submatches that are part of the
// supplied pattern, we return the first submatch. A zero value string
// is returned if nothing is found.
func fileContentMatch(path string, pattern string) (ret string) {
	fd, err := os.Open(path)
	if err != nil {
		return
	}
	defer func() {
		fd.Close()
	}()
	scanner := bufio.NewScanner(fd)

	re, err := regexp.Compile(pattern)
	if err != nil {
		return
	}

	for scanner.Scan() {
		buf := scanner.Text()
		subs := re.FindStringSubmatch(buf)
		if len(subs) >= 2 {
			ret = subs[1]
			break
		} else if len(subs) == 1 {
			ret = subs[0]
			break
		}
	}
	return
}
