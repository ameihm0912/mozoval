package oval

import (
	"os/exec"
	"strings"
)

type rpmRequest struct {
	out  chan rpmResponse
	name string
}

type rpmResponse struct {
	pkgdata rpmPackage
}

type rpmPackage struct {
	name    string
	version string
}

func (obj *GRPMInfoTest) execute(od *GOvalDefinitions) bool {
	v := od.getObject(obj.Object.ObjectRef)
	if v == nil {
		panic("unknown object in test execution!")
	}
	// XXX We should validate the object type here.
	o := v.(*GRPMInfoObj)

	s := od.getState(obj.State.StateRef)
	if s == nil {
		panic("unknown state in test execution")
	}
	// XXX We should validate the state type here.
	state := s.(*GRPMInfoState)

	return state.evaluate(o)
}

func (state *GRPMInfoState) evaluate(obj *GRPMInfoObj) bool {
	debugPrint("[rpminfo_state] evaluate %v\n", state.ID)

	transpkg := obj.Name
	if parserCfg.centosRedhatKludge != 0 {
		transpkg = centosRedhatPackageTranslate(transpkg)
	}

	rif := rpmRequest{}
	rif.out = make(chan rpmResponse)
	rif.name = transpkg
	dmgr.rpm.schan <- rif
	resp := <-rif.out

	// If we get nothing back the package isn't installed.
	if resp.pkgdata.name == "" {
		debugPrint("[rpminfo_state] doesn't look like %v is installed\n", transpkg)
		return false
	}
	debugPrint("[rpminfo_state] %v installed, %v\n", resp.pkgdata.name, resp.pkgdata.version)

	// If it's simply a key ID check, just simulate TRUE detection here.
	if len(state.SigKeyID.Value) > 0 {
		return true
	} else if len(state.EVRCheck.Value) > 0 {
		evrop := evrLookupOperation(state.EVRCheck.Operation)
		if evrop == EVROP_UNKNOWN {
			panic("evaluate: unknown evr comparison operation")
		}
		return evrCompare(evrop, resp.pkgdata.version, state.EVRCheck.Value)
	} else if len(state.VersionCheck.Value) > 0 {
		return versionPtrnMatch(resp.pkgdata.version, state.VersionCheck.Value)
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

	for _, x := range d.pkglist {
		if req.name == x.name {
			ret.pkgdata = x
			break
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
	buf, ok := c.Output()
	if ok != nil {
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
