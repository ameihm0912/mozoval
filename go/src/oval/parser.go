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
	flag_debug bool
}

var parser_cfg Config
var od GOvalDefinitions

func default_parser_config() Config {
	cfg := Config{
		flag_debug: false,
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

func Init() {
	parser_cfg = default_parser_config()
}

func Parse(path string) (*GOvalDefinitions, error) {
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
