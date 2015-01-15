package oval

func (od GOvalDefinitions) get_test(s string) interface{} {
	for _, x := range od.Tests.RPMInfoTests {
		if x.ID == s {
			return &x
		}
	}
	for _, x := range od.Tests.DPKGInfoTests {
		if x.ID == s {
			return &x
		}
	}

	return nil
}
