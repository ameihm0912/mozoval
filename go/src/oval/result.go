// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Contributor:
// - Aaron Meihm ameihm@mozilla.com
package oval

const (
	_ = iota
	RESULT_TRUE
	RESULT_FALSE
	RESULT_ERROR
)

//
// The result of an OVAL check
//
type GOvalResult struct {
	Status int
	Title  string
	ID     string
	Errors []string
}

func (gr *GOvalResult) StatusString() string {
	switch gr.Status {
	case RESULT_TRUE:
		return "true"
	case RESULT_FALSE:
		return "false"
	case RESULT_ERROR:
		return "error"
	}
	return "unknown"
}
