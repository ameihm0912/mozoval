package oval

import (
	"io/ioutil"
	"path/filepath"
)

func (obj *GTFC54Test) execute(od *GOvalDefinitions) bool {
	v := od.get_object(obj.Object.ObjectRef)
	if v == nil {
		// This should never happen as if the object doesnt exist we
		// would have seen that during preparation
		panic("unknown object in test execution!")
	}
	// XXX We should validate the object type here
	o := v.(*GTFC54Obj)

	s := od.get_state(obj.State.StateRef)
	if s == nil {
		panic("unknown state in test execution!")
	}
	state := s.(*GTFC54State)

	return state.evaluate(o)
}

func (obj *GTFC54Obj) resolve_path() (ret string) {
	if obj.Filepath != "" {
		ret = obj.Filepath
		return
	}

	l, err := ioutil.ReadDir(obj.Path)
	if err != nil {
		return
	}
	for _, x := range l {
		if x.Name() == obj.Filename {
			ret = filepath.Join(obj.Path, x.Name())
			break
		}
	}
	return
}

func (state *GTFC54State) evaluate(obj *GTFC54Obj) bool {
	debug_prt("[textfilecontent54_state] evaluate %v\n", state.ID)

	if obj.Pattern == "" {
		panic("textfilecontent54 evaluate with no pattern")
	}
	if state.SubExpression == "" {
		panic("textfilecontent54 evaluate with no subexpression")
	}

	path := obj.resolve_path()
	debug_prt("[textfilecontent54_state] target %v\n", path)
	debug_prt("[textfilecontent54_state] pattern %v\n", obj.Pattern)
	cmatch := file_content_match(path, obj.Pattern)
	if len(cmatch) == 0 {
		return false
	}
	debug_prt("[textfilecontent54_state] matched %v\n", cmatch)
	debug_prt("[textfilecontent54_state] compare %v\n",
		state.SubExpression)
	if state.SubExpression == cmatch {
		return true
	}

	return false
}

func (obj *GTFC54Obj) prepare() {
}
