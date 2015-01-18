package oval

import (
	"regexp"
	"strconv"
	"strings"
)

const (
	_ = iota
	EVROP_LESS_THAN
	EVROP_EQUALS
	EVROP_UNKNOWN
)

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

// Asset an epoch is present within a version string, if not a modified
// string is returned including a default epoch value (0)
func evr_epoch_assert(s string) string {
	f, _ := regexp.MatchString("^\\d+\\:", s)
	if !f {
		return "0:" + s
	}
	return s
}

func evr_extract(s string) (string, string, string) {
	var epoch string
	var version string
	var release string

	s0 := strings.Split(s, ":")
	if len(s0) < 2 {
		panic("evr_extract: can't extract epoch")
	}
	epoch = s0[0]

	// If we have a + character in the vr component, we treat this as a
	// dpkg style package, otherwise rpm
	if strings.Contains(s0[1], "+") {
		s0 = strings.Split(s0[1], "+")
		if len(s0) < 2 {
			panic("evr_extract: + tokenize failure")
		}
		version = s0[0]
		release = s0[1]
	} else {
		version = s0[1]
		release = ""
	}

	debug_prt("[evr_extract] epoch=%v, version=%v, revision=%v\n", epoch, version, release)
	return epoch, version, release
}

func evr_e_compare(actual string, check string) int {
	ai, err := strconv.Atoi(actual)
	if err != nil {
		panic("evr_e_compare: atoi actual")
	}
	ci, err := strconv.Atoi(check)
	if err != nil {
		panic("evr_e_compare: atoi actual")
	}
	if ai > ci {
		return 1
	} else if ai < ci {
		return -1
	}
	return 0
}

//
// Compare a component of a version string containing an integer followed
// by a character
//
func evr_v_compare_numalpha(actual string, check string) (int, bool) {
	abuf := strings.Split(actual, "")
	cbuf := strings.Split(check, "")
	for i, v := range cbuf {
		if i >= len(actual) {
			// The check version component has more elements then
			// the actual, return greater
			return 1, true
		}
		if v > abuf[i] {
			return -1, true
		} else if v < abuf[i] {
			return 1, true
		}
	}
	return 0, true
}

func evr_v_compare(actual string, check string) int {
	if len(actual) == 0 || len(check) == 0 {
		panic("evr_v_compare: empty version string")
	}
	debug_prt("[evr_v_compare] %v %v\n", actual, check)
	dashbuf_actual := strings.Split(actual, "-")
	dashbuf_check := strings.Split(check, "-")

	ret := 0
	for x, checkdash := range dashbuf_check {
		if x >= len(dashbuf_actual) {
			// The actual string has more dash components then the
			// comparison string does, return what we have so far
			// and ignore the rest
			return ret
		}
		// sigma represents the component of the dash buffer from the
		// actual value for this cycle
		sigma := dashbuf_actual[x]

		dot_check := strings.Split(checkdash, ".")
		dot_sig := strings.Split(sigma, ".")

		// Loop through each dot component in the version string;
		// regular integer values are handled simply, if the component
		// has other types of characters we pass them off to extended
		// handling functions
		for y, checkdot := range dot_check {
			if y >= len(dot_sig) {
				// There are more version components in this
				// string then in the check version, treat this
				// as greater if we have gotten this far
				return 1
			}
			ci, err_a := strconv.Atoi(checkdot)
			ai, err_c := strconv.Atoi(dot_sig[y])

			// If the conversion failed for either one, try a few
			// other extended comparison methods for the component
			extend := true
			if err_a != nil || err_c != nil {
				extend = false
				status, valid := evr_v_compare_numalpha(dot_sig[y], checkdot)
				if valid {
					extend = true
					if status > 0 {
						return 1
					} else if status < 0 {
						return -1
					}
				}
			}
			if !extend {
				panic("evr_v_compare: conversion and extended methods failed")
			}

			if ai > ci {
				return 1
			} else if ai < ci {
				return -1
			}
			// Otherwise the components were equal, continue on with the next
			// one
		}
	}
	return ret
}

func evr_r_compare(actual string, check string) int {
	return 0
}

func evr_compare(op int, actual string, check string) bool {
	debug_prt("[evr_compare] %v %v %v\n", actual, evr_operation_str(op), check)

	actual = evr_epoch_assert(actual)
	check = evr_epoch_assert(check)
	a_e, a_v, a_r := evr_extract(actual)
	c_e, c_v, c_r := evr_extract(check)

	res_epoch := evr_e_compare(a_e, c_e)
	res_version := evr_v_compare(a_v, c_v)
	res_release := evr_r_compare(a_r, c_r)
	debug_prt("[evr_compare] [%v:%v:%v] \n", res_epoch, res_version, res_release)

	switch op {
	case EVROP_EQUALS:
		if res_epoch == 0 &&
			res_version == 0 &&
			res_release == 0 {
			return true
		}
		return false
	case EVROP_LESS_THAN:
		switch res_epoch {
		case -1:
			return true
		case 1:
			return false
		}
		switch res_version {
		case -1:
			return true
		case 1:
			return false
		}
		switch res_release {
		case -1:
			return true
		case 1:
			return false
		}
		return false
	default:
		panic("unknown evr comparison operation")
	}
	return false
}

func Test_evr_compare(op int, actual string, check string) bool {
	return evr_compare(op, actual, check)
}
