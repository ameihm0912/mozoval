// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Contributor:
// - Aaron Meihm ameihm@mozilla.com
package oval

import (
	"io/ioutil"
	"path/filepath"
)

func (obj *GTFC54Test) execute(od *GOvalDefinitions, ctx defExecContext) (bool, error) {
	v := od.getObject(obj.Object.ObjectRef)
	if v == nil {
		ret := &ParserError{"unknown object in test execution"}
		ctx.error(ret.Error())
		return false, ret
	}

	o, ok := v.(*GTFC54Obj)
	if !ok {
		ret := &ParserError{"object is not textfilecontent54_object"}
		ctx.error(ret.Error())
		return false, ret
	}

	s := od.getState(obj.State.StateRef)
	if s == nil {
		ret := &ParserError{"unknown state in test execution"}
		ctx.error(ret.Error())
		return false, ret
	}

	state, ok := s.(*GTFC54State)
	if !ok {
		ret := &ParserError{"state is not textfilecontent54_state"}
		ctx.error(ret.Error())
		return false, ret
	}

	return state.evaluate(o), nil
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
	cmatch := fileContentMatchAll(path, obj.Pattern)
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
