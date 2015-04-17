// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Contributor:
// - Aaron Meihm ameihm@mozilla.com
package oval

import (
	"os/exec"
	"strings"
)

const (
	_ = iota
	DPKG_EXACT_MATCH
	DPKG_SUBSTRING_MATCH
)

type dpkgRequest struct {
	out       chan dpkgResponse
	name      string
	matchtype int
}

type dpkgResponse struct {
	pkgdata []dpkgPackage
}

type dpkgPackage struct {
	name    string
	version string
}

func (d *dpkgPackage) externalize() (ret ExternalizedPackage) {
	ret.Name = d.name
	ret.Version = d.version
	ret.PkgType = "dpkg"
	return ret
}

func (obj *GDPKGInfoTest) execute(od *GOvalDefinitions) bool {
	v := od.getObject(obj.Object.ObjectRef)
	if v == nil {
		panic("unknown object in test execution!")
	}
	// XXX We should validate the object type here.
	o := v.(*GDPKGInfoObj)

	s := od.getState(obj.State.StateRef)
	if s == nil {
		panic("unknown state in test execution")
	}
	// XXX We should validate the state type here.
	state := s.(*GDPKGInfoState)

	return state.evaluate(o)
}

func (state *GDPKGInfoState) evaluate(obj *GDPKGInfoObj) bool {
	debugPrint("[dpkginfo_state] evaluate %v\n", state.ID)

	resp := dmgr.dpkg.makeRequest(obj.Name, DPKG_EXACT_MATCH)

	// If we get nothing back the package isn't installed.
	if len(resp.pkgdata) == 0 {
		debugPrint("[dpkginfo_state] doesn't look like %v is installed\n", obj.Name)
		return false
	}

	// XXX It's possible multiple responses can be returned, right now we
	// just select the first one but we should probably sort and use the
	// latest.
	pkgname := resp.pkgdata[0].name
	pkgversion := resp.pkgdata[0].version
	debugPrint("[dpkginfo_state] %v installed, %v\n", pkgname, pkgversion)

	if len(state.EVRCheck.Value) > 0 {
		evrop := evrLookupOperation(state.EVRCheck.Operation)
		if evrop == EVROP_UNKNOWN {
			panic("evaluate: unknown evr comparison operation")
		}
		return evrCompare(evrop, pkgversion, state.EVRCheck.Value)
	}

	return false
}

func (d *GDPKGInfoObj) prepare() {
}

type dpkgDataMgr struct {
	schan    chan dpkgRequest
	pkglist  []dpkgPackage
	prepared bool
}

func (d *dpkgDataMgr) makeRequest(arg string, matchType int) dpkgResponse {
	if !d.prepared {
		panic("dpkg package manager not prepared")
	}

	dif := dpkgRequest{}
	dif.out = make(chan dpkgResponse)
	dif.name = arg
	dif.matchtype = matchType
	dmgr.dpkg.schan <- dif
	return <-dif.out
}

func (d *dpkgDataMgr) init() {
	debugPrint("initializing dpkg data manager\n")
	d.schan = make(chan dpkgRequest)
}

func (d *dpkgDataMgr) prepare() {
	d.pkglist = dpkgGetPackages()
	d.prepared = true
}

func (d *dpkgDataMgr) build_response(req dpkgRequest) dpkgResponse {
	ret := dpkgResponse{}

	for _, x := range d.pkglist {
		switch req.matchtype {
		case DPKG_EXACT_MATCH:
			if req.name == x.name {
				ret.pkgdata = append(ret.pkgdata, x)
			}
		case DPKG_SUBSTRING_MATCH:
			if strings.Contains(x.name, req.name) {
				ret.pkgdata = append(ret.pkgdata, x)
			}
		default:
			panic("invalid dpkg match type specified")
		}
	}

	return ret
}

func (d *dpkgDataMgr) run() {
	debugPrint("Starting dpkg data manager\n")

	for {
		r, ok := <-d.schan
		if ok == false {
			debugPrint("Stopping dpkg data manager\n")
			return
		}
		r.out <- d.build_response(r)
	}
}

func dpkgGetPackages() []dpkgPackage {
	ret := make([]dpkgPackage, 0)

	c := exec.Command("dpkg", "-l")
	buf, ok := c.Output()
	if ok != nil {
		return nil
	}

	slist := strings.Split(string(buf), "\n")
	for _, x := range slist {
		s := strings.Fields(x)

		if len(s) < 3 {
			continue
		}
		// Only process packages that have been fully installed.
		if s[0] != "ii" {
			continue
		}
		newpkg := dpkgPackage{s[1], s[2]}
		ret = append(ret, newpkg)
	}
	return ret
}
