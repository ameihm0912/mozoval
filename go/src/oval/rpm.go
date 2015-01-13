package oval

const (
	RPMCMD_GET = 0
)

type RPMRequest struct {
	command		int
}

type RPMResponse struct {
}

func handle_rpm_request(RPMRequest) RPMResponse {
	var ret RPMResponse
	return ret
}

func RPMManager(rch chan RPMRequest, wch chan RPMResponse) {
	for {
		req, ok := <- rch
		if ok == false {
			return
		}
		handle_rpm_request(req)
	}
}
