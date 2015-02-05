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

func (gc *GCriteria) statusString() string {
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

func (gc *GCriteria) evaluate(p *GOvalDefinitions) {
	// if the operator hasn't been set on the criteria default it to AND.
	// according to the OVAL spec the operator is a required field but it
	// seems a lot of definitions do not always include it.
	if gc.Operator == "" {
		gc.Operator = "AND"
	}

	if (gc.Operator != "AND") && (gc.Operator != "OR") {
		debugPrint("[criteria] criteria has invalid operator, ignoring\n")
		return
	}

	debugPrint("[criteria] %v\n", gc.Operator)

	// Evaluate all criteria, criterion, and extended definitions that are
	// part of this criteria element.
	for _, c := range gc.Criteria {
		c.evaluate(p)
	}
	for _, c := range gc.ExtendDef {
		c.evaluate(p)
	}
	for _, c := range gc.Criterion {
		c.evaluate(p)
	}

	debugPrint("[criteria] status=%v\n", gc.statusString())
}

func (gc *GCriterion) evaluate(p *GOvalDefinitions) {
	var tiface genericTest
	var result bool

	debugPrint("[criterion] %v\n", gc.Comment)

	r := p.getTest(gc.Test)
	if r == nil {
		debugPrint("[criterion] can't locate test %v\n", gc.Test)
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
		debugPrint("[criterion] unhandled test struct %v\n", reflect.TypeOf(r))
		gc.status = CRITERIA_ERROR
		return
	}

	tiface.prepare(p)

	result = tiface.execute(p)
	if result {
		debugPrint("[criterion] TRUE\n")
		gc.status = CRITERIA_TRUE
	} else {
		debugPrint("[criterion] FALSE\n")
		gc.status = CRITERIA_FALSE
	}
}

func (gc *GExtendDefinition) evaluate(p *GOvalDefinitions) {
	debugPrint("[extend_definition] %v\n", gc.Comment)

	x := p.getDefinition(gc.Test)
	if x == nil {
		debugPrint("can't locate definition %v\n", gc.Test)
		return
	}
	x.evaluate(nil, p)
}
