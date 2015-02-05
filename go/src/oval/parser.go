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
	flagDebug bool
	maxChecks int
}

type dataMgr struct {
	dpkg dpkgDataMgr
}

func (d *dataMgr) dataMgrInit() {
	d.dpkg.init()
}

func (d *dataMgr) dataMgrRun(precognition bool) {
	if precognition {
		d.dpkg.prepare()
	}
	go d.dpkg.run()
}

func (d *dataMgr) dataMgrClose() {
	close(d.dpkg.schan)
}

var parserCfg config
var dmgr dataMgr

func defaultParserConfig() config {
	return config{
		flagDebug: false,
		maxChecks: 10,
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
	fmt.Fprintf(os.Stderr, s, args...)
}

func Execute(od *GOvalDefinitions) {
	var precognition bool = false
	debugPrint("executing all applicable checks\n")

	if parserCfg.flagDebug {
		precognition = true
	}

	dmgr.dataMgrInit()
	dmgr.dataMgrRun(precognition)

	results := make([]GOvalResult, 0)
	reschan := make(chan GOvalResult)
	curchecks := 0
	for _, v := range od.Definitions.Definitions {
		debugPrint("executing definition %s...\n", v.ID)

		for {
			nodata := false
			select {
			case s := <-reschan:
				results = append(results, s)
				curchecks--
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
		}
		go v.evaluate(reschan, od)
		curchecks++
	}

	dmgr.dataMgrClose()
}

func Init() {
	parserCfg = defaultParserConfig()
}

func Parse(path string) (*GOvalDefinitions, error) {
	var od GOvalDefinitions
	var perr ParserError

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
