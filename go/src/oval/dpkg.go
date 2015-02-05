package oval

import (
	"os/exec"
	"strings"
)

type dpkgRequest struct {
	out  chan dpkgResponse
	name string
}

type dpkgResponse struct {
	pkgdata dpkgPackage
}

type dpkgPackage struct {
	name    string
	version string
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

	dif := dpkgRequest{}
	dif.out = make(chan dpkgResponse)
	dif.name = obj.Name
	dmgr.dpkg.schan <- dif
	resp := <-dif.out

	// If we get nothing back the package isn't installed.
	if resp.pkgdata.name == "" {
		debugPrint("[dpkginfo_state] doesn't look like %v is installed\n", obj.Name)
		return false
	}
	debugPrint("[dpkginfo_state] %v installed, %v\n", resp.pkgdata.name, resp.pkgdata.version)

	if len(state.EVRCheck.Value) > 0 {
		evrop := evrLookupOperation(state.EVRCheck.Operation)
		if evrop == EVROP_UNKNOWN {
			panic("evaluate: unknown evr comparison operation")
		}
		evrCompare(evrop, resp.pkgdata.version, state.EVRCheck.Value)
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
		if req.name == x.name {
			ret.pkgdata = x
			break
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
