package oval

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

const (
	_ = iota
	EVROP_LESS_THAN
	EVROP_EQUALS
	EVROP_UNKNOWN
)

type EVR struct {
	epoch   string
	version string
	release string
}

func evr_lookup_operation(s string) int {
	switch s {
	case "less than":
		return EVROP_LESS_THAN
	}
	return EVROP_UNKNOWN
}

func evr_operation_str(val int) string {
	switch val {
	case EVROP_LESS_THAN:
		return "<"
	case EVROP_EQUALS:
		return "="
	default:
		return "?"
	}
}

func evr_isdigit(c rune) bool {
	return unicode.IsDigit(c)
}

func evr_extract(s string) EVR {
	var ret EVR
	var idx int

	for _, c := range s {
		if !evr_isdigit(c) {
			break
		}
		idx++
	}

	if idx >= len(s) {
		panic("evr_extract: all digits")
	}

	if s[idx] == ':' {
		ret.epoch = s[:idx]
	} else {
		ret.epoch = "0"
	}

	idx++
	if idx >= len(s) {
		panic("evr_extract: only epoch")
	}
	remain := s[idx:]

	rp0 := strings.LastIndex(remain, "-")
	if rp0 != -1 {
		ret.version = remain[:rp0]
		rp0++
		if rp0 >= len(remain) {
			panic("evr_extract: ends in dash")
		}
		ret.release = remain[rp0:]
	} else {
		ret.version = remain
		ret.release = ""
	}

	debug_prt("[evr_extract] epoch=%v, version=%v, revision=%v\n",
		ret.epoch, ret.version, ret.release)
	return ret
}

func evr_rpmtokenizer(s string) []string {
	re := regexp.MustCompile("[A-Za-z0-9]+")
	return re.FindAllString(s, -1)
}

func evr_trimzeros(s string) string {
	if len(s) == 1 {
		return s
	}
	_, err := strconv.Atoi(s)
	if err != nil {
		return s
	}
	return strings.TrimLeft(s, "0")
}

func evr_rpmvercmp(actual string, check string) int {
	if actual == check {
		return 0
	}

	acttokens := evr_rpmtokenizer(actual)
	chktokens := evr_rpmtokenizer(check)

	for i := range chktokens {
		if i >= len(acttokens) {
			// There are more tokens in the check value, the
			// check wins
			return 1
		}

		// If the values are pure numbers, trim any leading 0's
		acttest := evr_trimzeros(acttokens[i])
		chktest := evr_trimzeros(chktokens[i])

		// Do a lexical string comparison here, this should work
		// even with pure integer values
		if chktest > acttest {
			return 1
		} else if chktest < acttest {
			return -1
		}
	}

	// If we get this far, see if the actual value still has more tokens
	// for comparison, if so actual wins
	if len(acttokens) > len(chktokens) {
		return -1
	}

	return 0
}

func evr_rpmcompare(actual EVR, check EVR) int {
	aepoch, err := strconv.Atoi(actual.epoch)
	if err != nil {
		panic("evr_rpmcompare: bad actual epoch")
	}
	cepoch, err := strconv.Atoi(check.epoch)
	if err != nil {
		panic("evr_rpmcompare: bad check epoch")
	}
	if cepoch > aepoch {
		return 1
	} else if cepoch < aepoch {
		return -1
	}

	ret := evr_rpmvercmp(actual.version, check.version)
	if ret != 0 {
		return ret
	}

	ret = evr_rpmvercmp(actual.release, check.release)
	if ret != 0 {
		return ret
	}

	return 0
}

func evr_compare(op int, actual string, check string) bool {
	debug_prt("[evr_compare] %v %v %v\n", actual, evr_operation_str(op),
		check)

	evract := evr_extract(actual)
	evrchk := evr_extract(check)

	ret := evr_rpmcompare(evract, evrchk)
	switch op {
	case EVROP_EQUALS:
		if ret != 0 {
			return false
		}
		return true
	case EVROP_LESS_THAN:
		if ret != 1 {
			return false
		}
		return true
	}
	panic("evr_compare: unknown operator")
}

func Test_evr_compare(op int, actual string, check string) bool {
	return evr_compare(op, actual, check)
}
