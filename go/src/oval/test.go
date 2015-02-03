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
	release()
	execute(*GOvalDefinitions) bool
}

func (gt *GTest) release() {
	gt.Unlock()
}

func (gt *GTest) prepare(od *GOvalDefinitions) {
	var iface genericobj

	gt.Lock()

	//
	// Prepare the object the test depends on, and return the state the
	// test applies to
	//
	v := od.get_object(gt.Object.ObjectRef)
	if v == nil {
		debug_prt("[test] can't locate object %s\n",
			gt.Object.ObjectRef)
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
	case reflect.TypeOf(&GTFC54Obj{}):
		r := v.(*GTFC54Obj)
		iface = r
	default:
		debug_prt("[test] unhandled object struct %v\n",
			reflect.TypeOf(v))
		gt.status = TEST_ERROR
		return
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
	for _, x := range od.Tests.TFC54Tests {
		if x.ID == s {
			return &x
		}
	}

	return nil
}
