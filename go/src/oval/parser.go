package oval

import (
	"fmt"
	"os"
	"encoding/xml"
)

type ParserError struct {
	s string
}
func (pe *ParserError) Error() string {
	return pe.s
}

type config struct {
	flag_debug	bool
	max_checks	int
}

type datamgr struct {
	dpkg		dpkgdatamgr
}

var parser_cfg config
var dmgr datamgr

func default_parser_config() config {
	cfg := config{
		flag_debug: false,
		// The maximum number of checks that can be run at any given
		// time, not configurable at the moment but should be
		max_checks: 10,
	}
	return cfg
}

func (d *datamgr) datamgr_init() {
	d.dpkg.init()
}

func (d *datamgr) datamgr_run(precognition bool) {
	if precognition {
		d.dpkg.prepare()
	}
	go d.dpkg.run()
}

func (d *datamgr) datamgr_close() {
	close(d.dpkg.schan)
}

func Set_debug(f bool) {
	parser_cfg.flag_debug = f
}

func debug_prt(s string, args ...interface{}) {
	if !parser_cfg.flag_debug {
		return
	}
	fmt.Fprintf(os.Stderr, s, args...)
}

func Execute(od *GOvalDefinitions) {
	var precognition bool = false
	debug_prt("Executing all applicable checks\n")

	if parser_cfg.flag_debug {
		precognition = true
	}

	dmgr.datamgr_init()
	dmgr.datamgr_run(precognition)

	results := make([]GOvalResult, 0)
	reschan := make(chan GOvalResult)
	curchecks := 0
	for _, v := range od.Definitions.Definitions {
		debug_prt("Executing definition %s...\n", v.ID)

		for {
			nodata := false
			select {
			case s := <- reschan:
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

		if curchecks == parser_cfg.max_checks {
			// Block and wait for a free slot
			s := <- reschan
			results = append(results, s)
			curchecks--
		}
		go v.evaluate(reschan)
		curchecks++
	}

	dmgr.datamgr_close()
}

func Init() {
	parser_cfg = default_parser_config()
}

func Parse(path string) (*GOvalDefinitions, error) {
	var od GOvalDefinitions
	var perr ParserError

	debug_prt("Parsing %s\n", path)

	xfd, err := os.Open(path)
	if err != nil {
		perr.s = fmt.Sprintf("Error opening file: %v", err)
		return nil, &perr
	}

	decoder := xml.NewDecoder(xfd)
	ok := decoder.Decode(&od)
	if ok != nil {
		perr.s = fmt.Sprintf("Error parsing %v: invalid XML format?", path)
		return nil, &perr
	}
	xfd.Close()

	return &od, nil
}
