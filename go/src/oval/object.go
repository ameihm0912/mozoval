package oval

type genericobj interface {
	prepare()
}

func (od *GOvalDefinitions) get_object(s string) interface{} {
	for _, x := range od.Objects.RPMInfoObjects {
		if x.ID == s {
			return &x
		}
	}
	for _, x := range od.Objects.DPKGInfoObjects {
		if x.ID == s {
			return &x
		}
	}
	for _, x := range od.Objects.TFC54Objects {
		if x.ID == s {
			return &x
		}
	}

	return nil
}
