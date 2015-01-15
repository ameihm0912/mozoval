package oval

func (od GOvalDefinitions) get_definition(s string) *GDefinition {
	for _, x := range od.Definitions.Definitions {
		if x.ID == s {
			return &x
		}
	}

	return nil
}

func (od GDefinition) evaluate(ch chan GOvalResult, p *GOvalDefinitions) {
	var ret GOvalResult

	od.Lock()

	debug_prt("[evaluate] %s\n", od.ID)

	// Evaluate the root criteria item, this could result in recursion
	// through subelements of the definition
	od.Criteria.evaluate(p)

	// If the channel was nil we don't send the result back, we only want
	// one result per definition
	if ch != nil {
		ch <- ret
	}

	od.Unlock()
}
