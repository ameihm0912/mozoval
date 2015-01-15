package oval

const (
	_ = iota
	CRITERIA_PASS
	CRITERIA_FAIL
	CRITERIA_ERROR
)

func (gc *GCriteria) status_string() string {
	switch gc.status {
	case CRITERIA_PASS:
		return "pass"
	case CRITERIA_FAIL:
		return "fail"
	case CRITERIA_ERROR:
		return "error"
	}
	return "unknown"
}

func (gc *GExtendDefinition) evaluate(p *GOvalDefinitions) {
	debug_prt("[extend_definition] %v\n", gc.Comment)

	x := p.get_definition(gc.Test)
	if x == nil {
		debug_prt("cant find definition %v\n", gc.Test)
		return
	}
	x.evaluate(nil, p)
}

func (gc *GCriteria) evaluate(p *GOvalDefinitions) {
	// if the operator hasn't been set on the criteria default it to AND.
	// according to the OVAL spec the operator is a required field but it
	// seems a lot of definitions do not always include it
	if gc.Operator == "" {
		gc.Operator = "AND"
	}

	if (gc.Operator != "AND") && (gc.Operator != "OR") {
		debug_prt("[criteria] criteria has invalid operator, ignoring\n")
		return
	}

	debug_prt("[criteria] %v\n", gc.Operator)

	for _, c := range gc.Criteria {
		c.evaluate(p)
	}
	for _, c := range gc.ExtendDef {
		c.evaluate(p)
	}
	for _, c := range gc.Criterion {
		c.evaluate(p)
	}

	debug_prt("[criteria] status=%v\n", gc.status_string())
}

func (gc *GCriterion) evaluate(p *GOvalDefinitions) {
	debug_prt("[criterion] %v\n", gc.Comment)
}
