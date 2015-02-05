package oval

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

//
// Given a file, read the file line by line matching against pattern; if
// we find a match, return it. If there are submatches are part of the
// supplied pattern, we return the first submatch.
//
func fileContentMatch(path string, pattern string) (ret string) {
	var lastmatch = false

	fd, err := os.Open(path)
	if err != nil {
		return
	}
	rdr := bufio.NewReader(fd)
	re, err := regexp.Compile(pattern)
	if err != nil {
		fd.Close()
		return
	}
	for {
		buf, err := rdr.ReadString('\n')
		if err != nil {
			lastmatch = true
		}

		if len(buf) == 0 {
			return
		}

		subs := re.FindStringSubmatch(strings.Trim(buf, "\n"))
		if len(subs) >= 2 {
			ret = subs[1]
			break
		} else if len(subs) == 1 {
			ret = subs[0]
			break
		}
		if lastmatch {
			break
		}
	}
	fd.Close()
	return ret
}
