package oval

func (obj *GTFC54Test) execute(od *GOvalDefinitions) bool {
	v := od.get_object(obj.Object.ObjectRef)
	if v == nil {
		// This should never happen as if the object doesnt exist we
		// would have seen that during preparation
		panic("unknown object in test execution!")
	}
	// XXX We should validate the object type here
	_ = v.(*GTFC54Obj)

	s := od.get_state(obj.State.StateRef)
	if s == nil {
		panic("unknown state in test execution!")
	}

	return false
}

func (obj *GTFC54Obj) prepare() {
}
