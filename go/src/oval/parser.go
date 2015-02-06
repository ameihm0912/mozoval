// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Contributor:
// - Aaron Meihm ameihm@mozilla.com
package oval

import (
	"encoding/xml"
	"fmt"
	"os"
)

type ParserError struct {
	s string
}

func (pe *ParserError) Error() string {
	return pe.s
}

type config struct {
	flagDebug          bool
	maxChecks          int
	centosRedhatKludge int
}

type dataMgr struct {
	dpkg dpkgDataMgr
	rpm  rpmDataMgr
}

func (d *dataMgr) dataMgrInit() {
	d.dpkg.init()
	d.rpm.init()
}

func (d *dataMgr) dataMgrRun(precognition bool) {
	if precognition {
		d.dpkg.prepare()
		d.rpm.prepare()
	}
	go d.dpkg.run()
	go d.rpm.run()
}

func (d *dataMgr) dataMgrClose() {
	close(d.dpkg.schan)
	close(d.rpm.schan)
}

var parserCfg config
var dmgr dataMgr

func defaultParserConfig() config {
	return config{
		flagDebug:          false,
		maxChecks:          10,
		centosRedhatKludge: 0,
	}
}

func SetDebug(f bool) {
	parserCfg.flagDebug = f
}

func SetMaxChecks(i int) {
	parserCfg.maxChecks = i
}

func debugPrint(s string, args ...interface{}) {
	if !parserCfg.flagDebug {
		return
	}
	fmt.Fprintf(os.Stdout, s, args...)
}

func Execute(od *GOvalDefinitions) []GOvalResult {
	var precognition bool = false
	debugPrint("executing all applicable checks\n")

	precognition = true

	dmgr.dataMgrInit()
	dmgr.dataMgrRun(precognition)

	results := make([]GOvalResult, 0)
	reschan := make(chan GOvalResult)
	curchecks := 0
	expect := len(od.Definitions.Definitions)
	for _, v := range od.Definitions.Definitions {
		debugPrint("executing definition %s...\n", v.ID)

		for {
			nodata := false
			select {
			case s := <-reschan:
				results = append(results, s)
				curchecks--
				expect--
			default:
				nodata = true
				break
			}
			if nodata {
				break
			}
		}

		if curchecks == parserCfg.maxChecks {
			// Block and wait for a free slot.
			s := <-reschan
			results = append(results, s)
			curchecks--
			expect--
		}
		go v.evaluate(reschan, od)
		curchecks++
	}

	for expect > 0 {
		s := <-reschan
		results = append(results, s)
		expect--
	}

	dmgr.dataMgrClose()

	return results
}

func Init() {
	parserCfg = defaultParserConfig()
}

func Parse(path string) (*GOvalDefinitions, error) {
	var od GOvalDefinitions
	var perr ParserError

	parserCfg.centosRedhatKludge = centosDetection()

	debugPrint("parsing %s\n", path)

	xfd, err := os.Open(path)
	if err != nil {
		perr.s = fmt.Sprintf("error opening file: %v", err)
		return nil, &perr
	}

	decoder := xml.NewDecoder(xfd)
	ok := decoder.Decode(&od)
	if ok != nil {
		perr.s = fmt.Sprintf("error parsing %v: invalid xml format?", path)
		return nil, &perr
	}
	xfd.Close()

	return &od, nil
}
