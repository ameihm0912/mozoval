package oval

import (
	"os/exec"
	"strings"
)

type dpkgrequest struct {
	out  chan dpkgresponse
	name string
}

type dpkgresponse struct {
	pkgdata dpkgpackage
}

type dpkgdatamgr struct {
	schan    chan dpkgrequest
	pkglist  []dpkgpackage
	prepared bool
}

type dpkgpackage struct {
	name    string
	version string
}

func (obj *GDPKGInfoTest) execute(od *GOvalDefinitions) bool {
	v := od.get_object(obj.Object.ObjectRef)
	if v == nil {
		// This should never happen as if the object doesnt exist we
		// would have seen that during preparation
		panic("unknown object in test execution!")
	}
	// XXX We should validate the object type here
	o := v.(*GDPKGInfoObj)

	s := od.get_state(obj.State.StateRef)
	if s == nil {
		panic("unknown state in test execution!")
	}
	state := s.(*GDPKGInfoState)

	return state.evaluate(o)
}

func (state *GDPKGInfoState) evaluate(obj *GDPKGInfoObj) bool {
	debug_prt("[dpkginfo_state] evaluate %v\n", state.ID)

	dif := dpkgrequest{}
	dif.out = make(chan dpkgresponse)
	dif.name = obj.Name
	dmgr.dpkg.schan <- dif
	resp := <-dif.out

	// If we get nothing back the package isn't installed
	if resp.pkgdata.name == "" {
		debug_prt("[dpkginfo_state] doesn't look like %v is installed\n", obj.Name)
		return false
	}
	debug_prt("[dpkginfo_state] %v installed, %v\n", resp.pkgdata.name, resp.pkgdata.version)

	if len(state.EVRCheck.Value) > 0 {
		evrop := evr_lookup_operation(state.EVRCheck.Operation)
		if evrop == EVROP_UNKNOWN {
			return false
		}
		evr_compare(evrop, resp.pkgdata.version, state.EVRCheck.Value)
	}

	return false
}

func (d *GDPKGInfoObj) prepare() {
}

func (d *dpkgdatamgr) init() {
	debug_prt("Initializing dpkg data manager\n")
	d.schan = make(chan dpkgrequest)
}

func (d *dpkgdatamgr) prepare() {
	d.pkglist = dpkg_get_packages()
	d.prepared = true
}

func (d *dpkgdatamgr) build_response(req dpkgrequest) dpkgresponse {
	ret := dpkgresponse{}

	for _, x := range d.pkglist {
		if req.name == x.name {
			ret.pkgdata = x
			break
		}
	}

	return ret
}

func (d *dpkgdatamgr) run() {
	debug_prt("Starting dpkg data manager\n")

	for {
		r, ok := <-d.schan
		if ok == false {
			debug_prt("Stopping dpkg data manager\n")
			return
		}
		r.out <- d.build_response(r)
	}
}

func dpkg_get_packages() []dpkgpackage {
	ret := make([]dpkgpackage, 0)

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
		// If the package isn't installed ignore it
		if s[0] != "ii" {
			continue
		}
		newpkg := dpkgpackage{s[1], s[2]}
		ret = append(ret, newpkg)
	}
	return ret
}
