package oval

func (od GOvalDefinitions) getDefinition(s string) *GDefinition {
	for _, x := range od.Definitions.Definitions {
		if x.ID == s {
			return &x
		}
	}

	return nil
}

func (od *GOvalDefinitions) getState(s string) interface{} {
	for _, x := range od.States.RPMInfoStates {
		if x.ID == s {
			return &x
		}
	}
	for _, x := range od.States.TFC54States {
		if x.ID == s {
			return &x
		}
	}
	for _, x := range od.States.DPKGInfoStates {
		if x.ID == s {
			return &x
		}
	}

	return nil
}

func (od *GOvalDefinitions) getObject(s string) interface{} {
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

func (od *GOvalDefinitions) getTest(s string) interface{} {
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
	for _, x := range od.Tests.TFC54Tests {
		if x.ID == s {
			return &x
		}
	}

	return nil
}

func (od GDefinition) evaluate(ch chan GOvalResult, p *GOvalDefinitions) {
	var ret GOvalResult

	// We need a lock here as this definition could be selected for
	// evaluation by another definition as part of an extended
	// definition call.
	od.Lock()

	debugPrint("[evaluate] %v\n", od.ID)

	// Evaluate the root criteria item; this will likely result in
	// recursion through various subelements in the definition.
	od.status = od.Criteria.evaluate(p)
	ret.Status = od.status
	ret.Title = od.Metadata.Title
	ret.ID = od.ID

	// If the channel was nil we don't send the result back. This can
	// occur if the definition was called as the result of an
	// extend_definition rule in the OVAL definition being evaluated.
	if ch != nil {
		ch <- ret
	}

	od.Unlock()
}
