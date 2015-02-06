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
