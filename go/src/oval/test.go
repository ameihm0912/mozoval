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

type genericTest interface {
	prepare(*GOvalDefinitions)
	release()
	execute(*GOvalDefinitions) bool
}

func (gt *GTest) release() {
	gt.Unlock()
}

func (gt *GTest) prepare(od *GOvalDefinitions) {
	var iface genericObj

	gt.Lock()

	// Prepare the object the test depends on.
	v := od.getObject(gt.Object.ObjectRef)
	if v == nil {
		debugPrint("[test] can't locate object %s\n", gt.Object.ObjectRef)
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
		debugPrint("[test] unhandled object struct %v\n", reflect.TypeOf(v))
		gt.status = TEST_ERROR
		return
	}
	iface.prepare()
}
