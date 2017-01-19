package api

type simpleMessage struct {
	Body string `json:"body"`
}

type errorMessage struct {
	Error string `json:"error"`
}
