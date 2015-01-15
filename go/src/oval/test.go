package oval

import (
	"reflect"
)

const (
	_ = iota
	TEST_PASS
	TEST_FAIL
	TEST_ERROR
)

type generictest interface {
	prepare(*GOvalDefinitions)
}

func (gt *GTest) prepare(od *GOvalDefinitions) {
	var iface genericobj

	//
	// Prepare the object the test depends on
	//
	v := od.get_object(gt.Object.ObjectRef)
	if v == nil {
		debug_prt("[test] can't locate object %s\n", gt.Object.ObjectRef)
		gt.status = TEST_ERROR
		return
	}
	switch reflect.TypeOf(v) {
	case reflect.TypeOf(&GRPMInfoObj{}):
		r := v.(*GRPMInfoObj)
		iface = r
	case reflect.TypeOf(&GDPKGInfoObj{}):
		r := v.(*GDPKGInfoObj)
		iface = r
	default:
	}
	iface.prepare()
}

func (od *GOvalDefinitions) get_test(s string) interface{} {
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
