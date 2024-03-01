package types

type issue struct {
	Error string `json:"error"`
}

func NewIssue(errMsg string) issue {
	return issue{Error: errMsg}
}
