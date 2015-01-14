package oval

type DPKGRequest struct {
	out		chan DPKGResponse
}

type DPKGResponse struct {
}

type DPKGDataMgr struct {
	schan		chan DPKGRequest
	pkglist		[]DPKGPackage
}

type DPKGPackage struct {
	name		string
	version		string
}

func (d *DPKGDataMgr) init() {
	debug_prt("Initializing dpkg data manager\n")
	d.schan = make(chan DPKGRequest)
}

func (d *DPKGDataMgr) prepare() {
	d.pkglist = dpkg_get_packages()
}

func (d *DPKGDataMgr) run() {
	debug_prt("Starting dpkg data manager\n")

	for {
		_, ok := <- d.schan
		if ok == false {
			return
		}
	}
}

func dpkg_get_packages() []DPKGPackage {
	return nil
}
