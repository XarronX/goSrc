package types

type response struct {
	Message string `json:"message"`
}

func NewResponse(respMsg string) response {
	return response{Message: respMsg}
}
