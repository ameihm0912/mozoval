// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Contributor:
// - Aaron Meihm ameihm@mozilla.com
package oval

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
	iface = od.getObject(gt.Object.ObjectRef)
	if iface == nil {
		debugPrint("[test] can't locate object %s\n", gt.Object.ObjectRef)
		gt.status = TEST_ERROR
		return
	}
	iface.prepare()
}
