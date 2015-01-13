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

var parser_cfg Config

func default_parser_config() Config {
	cfg := Config{
		flag_debug: false,
		max_checks: 10,
	}
	return cfg
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
