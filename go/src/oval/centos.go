// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Contributor:
// - Aaron Meihm ameihm@mozilla.com
package oval

import (
	"regexp"
)

func centosRedhatPackageTranslate6(s string) string {
	ptrns := map[string]string{
		"redhat-release-server": "centos-release",
	}
	val, ok := ptrns[s]
	if !ok {
		return s
	}
	debugPrint("centos 6 translate: %v -> %v\n", s, val)
	return val
}

func centosRedhatPackageTranslate(s string) string {
	switch parserCfg.centosRedhatKludge {
	case 6:
		return centosRedhatPackageTranslate6(s)
	}
	return s
}

func centosDetection() int {
	centosVersionPatterns := map[string]int{
		"CentOS release 6\\..*": 6,
	}

	debugPrint("detecting centos\n")
	val := fileContentMatch("/etc/centos-release", "CentOS.*")
	if len(val) == 0 {
		debugPrint("centos not found\n")
		return 0
	}
	for k, v := range centosVersionPatterns {
		res, _ := regexp.MatchString(k, val)
		if res {
			debugPrint("detected centos %v\n", v)
			return v
		}
	}
	return 0
}
