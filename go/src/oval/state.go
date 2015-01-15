package oval

func (od *GOvalDefinitions) get_state(s string) interface{} {
	for _, x := range od.States.RPMInfoStates {
		if x.ID == s {
			return &x
		}
	}

	return nil
}
