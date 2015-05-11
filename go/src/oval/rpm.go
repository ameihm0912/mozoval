// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Contributor:
// - Aaron Meihm ameihm@mozilla.com
package oval

import (
	"os/exec"
	"regexp"
	"strings"
)

const (
	_ = iota
	RPM_EXACT_MATCH
	RPM_SUBSTRING_MATCH
	RPM_REGEXP_MATCH
)

type rpmRequest struct {
	out       chan rpmResponse
	name      string
	matchtype int
}

type rpmResponse struct {
	pkgdata []rpmPackage
}

type rpmPackage struct {
	name    string
	version string
}

func (r *rpmPackage) externalize() (ret ExternalizedPackage) {
	ret.Name = r.name
	ret.Version = r.version
	ret.PkgType = "rpm"
	return ret
}

func (obj *GRPMInfoTest) execute(od *GOvalDefinitions, ctx defExecContext) (bool, error) {
	v := od.getObject(obj.Object.ObjectRef)
	if v == nil {
		ret := &ParserError{"unknown object in test execution"}
		ctx.error(ret.Error())
		return false, ret
	}

	o, ok := v.(*GRPMInfoObj)
	if !ok {
		ret := &ParserError{"object is not rpminfo_object"}
		ctx.error(ret.Error())
		return false, ret
	}

	s := od.getState(obj.State.StateRef)
	if s == nil {
		ret := &ParserError{"unknown state in test execution"}
		ctx.error(ret.Error())
		return false, ret
	}

	state, ok := s.(*GRPMInfoState)
	if !ok {
		ret := &ParserError{"state is not rpminfo_state"}
		ctx.error(ret.Error())
		return false, ret
	}

	return state.evaluate(o), nil
}

func (state *GRPMInfoState) evaluate(obj *GRPMInfoObj) bool {
	debugPrint("[rpminfo_state] evaluate %v\n", state.ID)

	transpkg := obj.Name
	if parserCfg.centosRedhatKludge != 0 {
		transpkg = centosRedhatPackageTranslate(transpkg)
	}

	resp := dmgr.rpm.makeRequest(transpkg, RPM_EXACT_MATCH)

	// If we get nothing back the package isn't installed.
	if len(resp.pkgdata) == 0 {
		debugPrint("[rpminfo_state] doesn't look like %v is installed\n", transpkg)
		return false
	}

	// XXX It's possible multiple responses can be returned, right now we
	// just select the first one but we should probably sort and use the
	// latest.
	pkgname := resp.pkgdata[0].name
	pkgversion := resp.pkgdata[0].version
	debugPrint("[rpminfo_state] %v installed, %v\n", pkgname, pkgversion)

	// If it's simply a key ID check, just simulate TRUE detection here.
	if len(state.SigKeyID.Value) > 0 {
		return true
	} else if len(state.EVRCheck.Value) > 0 {
		evrop := evrLookupOperation(state.EVRCheck.Operation)
		if evrop == EVROP_UNKNOWN {
			panic("evaluate: unknown evr comparison operation")
		}
		return evrCompare(evrop, pkgversion, state.EVRCheck.Value)
	} else if len(state.VersionCheck.Value) > 0 {
		return versionPtrnMatch(pkgversion, state.VersionCheck.Value)
	}

	return false
}

func (r *GRPMInfoObj) prepare() {
}

type rpmDataMgr struct {
	schan    chan rpmRequest
	pkglist  []rpmPackage
	prepared bool
}

func (d *rpmDataMgr) makeRequest(arg string, matchType int) rpmResponse {
	if !d.prepared {
		panic("rpm package manager not prepared")
	}

	rif := rpmRequest{}
	rif.out = make(chan rpmResponse)
	rif.name = arg
	rif.matchtype = matchType
	dmgr.rpm.schan <- rif
	return <-rif.out
}

func (d *rpmDataMgr) init() {
	debugPrint("initializing rpm data manager\n")
	d.schan = make(chan rpmRequest)
}

func (d *rpmDataMgr) prepare() {
	d.pkglist = rpmGetPackages()
	d.prepared = true
}

func (d *rpmDataMgr) build_response(req rpmRequest) rpmResponse {
	ret := rpmResponse{}

	var rematch *regexp.Regexp
	var err error
	if req.matchtype == RPM_REGEXP_MATCH {
		rematch, err = regexp.Compile(req.name)
		if err != nil {
			return ret
		}
	}

	for _, x := range d.pkglist {
		switch req.matchtype {
		case RPM_EXACT_MATCH:
			if req.name == x.name {
				ret.pkgdata = append(ret.pkgdata, x)
			}
		case RPM_SUBSTRING_MATCH:
			if strings.Contains(x.name, req.name) {
				ret.pkgdata = append(ret.pkgdata, x)
			}
		case RPM_REGEXP_MATCH:
			if rematch.MatchString(x.name) {
				ret.pkgdata = append(ret.pkgdata, x)
			}
		default:
			panic("invalid rpm match type specified")
		}
	}

	return ret
}

func (d *rpmDataMgr) run() {
	debugPrint("Starting rpm data manager\n")

	for {
		r, ok := <-d.schan
		if ok == false {
			debugPrint("Stopping rpm data manager\n")
			return
		}
		r.out <- d.build_response(r)
	}
}

func rpmGetPackages() []rpmPackage {
	ret := make([]rpmPackage, 0)

	c := exec.Command("rpm", "-qa", "--queryformat", "%{NAME} %{EVR}\\n")
	buf, err := c.Output()
	if err != nil {
		return nil
	}

	slist := strings.Split(string(buf), "\n")
	for _, x := range slist {
		s := strings.Fields(x)

		if len(s) < 2 {
			continue
		}
		newpkg := rpmPackage{s[0], s[1]}
		ret = append(ret, newpkg)
	}
	return ret
}
