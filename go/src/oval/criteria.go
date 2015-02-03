package oval

import (
	"reflect"
)

const (
	_ = iota
	CRITERIA_TRUE
	CRITERIA_FALSE
	CRITERIA_ERROR
)

func (gc *GCriteria) status_string() string {
	switch gc.status {
	case CRITERIA_TRUE:
		return "true"
	case CRITERIA_FALSE:
		return "false"
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
		debug_prt("[criteria] criteria has invalid operator, " +
			"ignoring\n")
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
	var tiface generictest
	var result bool

	debug_prt("[criterion] %v\n", gc.Comment)

	r := p.get_test(gc.Test)
	if r == nil {
		debug_prt("[criterion] can't locate test %s\n", gc.Test)
		gc.status = CRITERIA_ERROR
		return
	}
	switch reflect.TypeOf(r) {
	case reflect.TypeOf(&GRPMInfoTest{}):
		v := r.(*GRPMInfoTest)
		tiface = v
	case reflect.TypeOf(&GDPKGInfoTest{}):
		v := r.(*GDPKGInfoTest)
		tiface = v
	case reflect.TypeOf(&GTFC54Test{}):
		v := r.(*GTFC54Test)
		tiface = v
	default:
		debug_prt("[criterion] unhandled test struct %v\n",
			reflect.TypeOf(r))
		gc.status = CRITERIA_ERROR
		return
	}

	tiface.prepare(p)

	result = tiface.execute(p)
	if result {
		debug_prt("[criterion] TRUE\n")
		gc.status = CRITERIA_TRUE
	} else {
		debug_prt("[criterion] FALSE\n")
		gc.status = CRITERIA_FALSE
	}
}
