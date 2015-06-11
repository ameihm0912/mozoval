// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Contributor:
// - Aaron Meihm ameihm@mozilla.com
package main

// Convert the output of the MIG pkg module OVAL interpreter into JSON events
// for consumption by MozDef

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ameihm0912/govfeed/src/govfeed"
	"github.com/ameihm0912/gozdef"
	"os"
	"regexp"
)

const sourceName = "mozoval"

var cveCache map[string]govfeed.GVCVE

var useVFeed string
var vulnEvents []gozdef.VulnEvent

func getCVE(cve string) (govfeed.GVCVE, error) {
	x, ok := cveCache[cve]
	if ok {
		return x, nil
	}
	x, err := govfeed.GVQuery(cve)
	if err != nil {
		return x, err
	}
	cveCache[cve] = x
	return x, nil
}

func lineParser(buf string) error {
	var (
		Hostname string
		CVE      string
		ID       string
		Title    string
	)

	var pTable = []struct {
		expression string
		target     *string
	}{
		{".*(CVE-\\d+-\\d+).*", &CVE},
		{".*id=(\\S+).*", &ID},
		{".*title=\"([^\"]+).*", &Title},
		{"^(\\S+) ovalresult.*", &Hostname},
	}

	for i := range pTable {
		r := regexp.MustCompile(pTable[i].expression)
		matches := r.FindStringSubmatch(buf)
		if len(matches) <= 1 {
			continue
		}
		*pTable[i].target = matches[1]
	}

	e, err := gozdef.NewVulnEvent()
	if err != nil {
		return err
	}
	e.SourceName = sourceName
	e.Description = fmt.Sprintf("mozoval check for %v", Hostname)

	// XXX These need to be set correctly
	e.Asset.AssetID = 1
	e.Vuln.VulnID = "mozoval-vuln"

	e.Asset.Hostname = Hostname

	if CVE != "" {
		e.Vuln.Title = CVE
		e.Vuln.CVE = append(e.Vuln.CVE, CVE)
	}

	if useVFeed != "" && CVE != "" {
		cvedata, err := getCVE(CVE)
		if err != nil {
			return err
		}
		if cvedata.Description != "" {
			e.Vuln.Description = cvedata.Description
		}
		if cvedata.CVSS != 0 {
			e.Vuln.CVSS = cvedata.CVSS
		}
	}

	if err = e.Validate(); err != nil {
		return err
	}

	vulnEvents = append(vulnEvents, e)
	return nil
}

func main() {
	flag.StringVar(&useVFeed, "v", "", "path to vFeed CLI")
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprint(os.Stderr, "error: specify results output file as argument\n")
		os.Exit(1)
	}

	if useVFeed != "" {
		err := govfeed.GVInit(useVFeed)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		cveCache = make(map[string]govfeed.GVCVE)
	}

	fd, err := os.Open(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	defer fd.Close()
	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		buf := scanner.Text()
		err = lineParser(buf)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}
	if err = scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	for _, x := range vulnEvents {
		jb, err := json.Marshal(&x)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "%v\n", string(jb))
	}
}
