package oval

func (od GDefinition) evaluate(ch chan GOvalResult) {
	var ret GOvalResult
	debug_prt("[evaluate] %s\n", od.ID)

	ch <- ret
}
