package oval

import (
	"os/exec"
	"strings"
)

type dpkgrequest struct {
	out		chan dpkgresponse
}

type dpkgresponse struct {
}

type dpkgdatamgr struct {
	schan		chan dpkgrequest
	pkglist		[]dpkgpackage
	prepared	bool
}

type dpkgpackage struct {
	name		string
	version		string
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

func (d *dpkgdatamgr) run() {
	debug_prt("Starting dpkg data manager\n")

	for {
		_, ok := <- d.schan
		if ok == false {
			debug_prt("Stopping dpkg data manager\n")
			return
		}
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
