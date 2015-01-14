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

type Config struct {
	flag_debug	bool
	max_checks	int
}

type DataMgr struct {
	dpkg		DPKGDataMgr
}

var parser_cfg Config
var datamgr DataMgr

func default_parser_config() Config {
	cfg := Config{
		flag_debug: false,
		// The maximum number of checks that can be run at any given
		// time, not configurable at the moment but should be
		max_checks: 10,
	}
	return cfg
}

func (d *DataMgr) datamgr_init() {
	d.dpkg.init()
}

func (d *DataMgr) datamgr_run() {
	go d.dpkg.run()
}

func (d *DataMgr) datamgr_close() {
	close(d.dpkg.schan)
}

func (d *DataMgr) precognition() {
	d.dpkg.prepare()
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
	debug_prt("Executing all applicable checks\n")

	datamgr.datamgr_init()
	datamgr.datamgr_run()
	if parser_cfg.flag_debug {
		datamgr.precognition()
	}

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
		go v.Evaluate(reschan)
		curchecks++
	}

	datamgr.datamgr_close()
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
