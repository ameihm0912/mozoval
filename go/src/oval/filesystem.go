package oval

import (
	"io/ioutil"
	"path/filepath"
)

func (obj *GTFC54Test) execute(od *GOvalDefinitions) bool {
	v := od.getObject(obj.Object.ObjectRef)
	if v == nil {
		panic("unknown object in test execution")
	}
	// XXX We should validate the object type here.
	o := v.(*GTFC54Obj)

	s := od.getState(obj.State.StateRef)
	if s == nil {
		panic("unknown state in test execution")
	}
	// XXX We should validate the state type here.
	state := s.(*GTFC54State)

	return state.evaluate(o)
}

func (obj *GTFC54Obj) resolvePath() (ret string) {
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

func (obj *GTFC54Obj) prepare() {
}

func (state *GTFC54State) evaluate(obj *GTFC54Obj) bool {
	debugPrint("[textfilecontent54_state] evaluate %v\n", state.ID)

	if obj.Pattern == "" {
		panic("textfilecontent54 evaluate with no pattern")
	}
	if state.SubExpression == "" {
		panic("textfilecontent54 evaluate with no subexpression")
	}

	path := obj.resolvePath()
	debugPrint("[textfilecontent54_state] target %v\n", path)
	debugPrint("[textfilecontent54_state] pattern %v\n", obj.Pattern)
	cmatch := fileContentMatch(path, obj.Pattern)
	if len(cmatch) == 0 {
		return false
	}
	debugPrint("[textfilecontent54_state] matched %v\n", cmatch)
	debugPrint("[textfilecontent54_state] compare %v\n",
		state.SubExpression)
	if state.SubExpression == cmatch {
		return true
	}

	return false
}
