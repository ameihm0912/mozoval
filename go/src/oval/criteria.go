// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Contributor:
// - Aaron Meihm ameihm@mozilla.com
package oval

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

func (gc *GCriteria) evaluate(p *GOvalDefinitions, ctx defExecContext) int {
	// if the operator hasn't been set on the criteria default it to AND.
	// according to the OVAL spec the operator is a required field but it
	// seems a lot of definitions do not always include it.
	if gc.Operator == "" {
		gc.Operator = "AND"
	}

	if (gc.Operator != "AND") && (gc.Operator != "OR") {
		ctx.error("[criteria] criteria has invalid operator, ignoring")
		gc.status = CRITERIA_ERROR
		return gc.status
	}

	debugPrint("[criteria] %v\n", gc.Operator)

	// Evaluate all criteria, criterion, and extended definitions that are
	// part of this criteria element.
	results := make([]int, 0)
	for i := range gc.Criteria {
		results = append(results, gc.Criteria[i].evaluate(p, ctx))
	}
	for i := range gc.ExtendDef {
		results = append(results, gc.ExtendDef[i].evaluate(p))
	}
	for i := range gc.Criterion {
		results = append(results, gc.Criterion[i].evaluate(p))
	}

	// If an error occurred anywhere during evaluation, return an error
	// for the criteria.
	for _, c := range results {
		if c == CRITERIA_ERROR {
			gc.status = CRITERIA_ERROR
			debugPrint("[criteria] status=%v\n", gc.statusString())
			return gc.status
		}
	}

	if gc.Operator == "AND" {
		gc.status = CRITERIA_TRUE
		for _, c := range results {
			if c != CRITERIA_TRUE {
				gc.status = CRITERIA_FALSE
				break
			}
		}
	} else {
		gc.status = CRITERIA_FALSE
		for _, c := range results {
			if c == CRITERIA_TRUE {
				gc.status = CRITERIA_TRUE
				break
			}
		}
	}
	debugPrint("[criteria] status=%v\n", gc.statusString())
	return gc.status
}

func (gc *GCriterion) evaluate(p *GOvalDefinitions) int {
	var tiface genericTest
	var result bool

	debugPrint("[criterion] %v\n", gc.Comment)

	tiface = p.getTest(gc.Test)
	if tiface == nil {
		debugPrint("[criterion] can't locate test %v\n", gc.Test)
		gc.status = CRITERIA_ERROR
		return gc.status
	}

	tiface.prepare(p)
	defer func() {
		tiface.release()
	}()

	result = tiface.execute(p)
	if result {
		debugPrint("[criterion] TRUE\n")
		gc.status = CRITERIA_TRUE
	} else {
		debugPrint("[criterion] FALSE\n")
		gc.status = CRITERIA_FALSE
	}
	return gc.status
}

func (gc *GExtendDefinition) evaluate(p *GOvalDefinitions) int {
	debugPrint("[extend_definition] %v\n", gc.Comment)

	x := p.getDefinition(gc.Test)
	if x == nil {
		debugPrint("can't locate definition %v\n", gc.Test)
		return CRITERIA_ERROR
	}
	x.evaluate(nil, p)
	return x.status
}
