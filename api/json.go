package api

// MsgResponse is an HTTP response in case of success
type MsgResponse struct {
	Msg string `json:"msg"`
}

// HealthResponse is an HTTP response for the 'health' endpoint
type HealthResponse struct {
	Health string `json:"health"`
}

// MsgRemoteResponse is an HTTP response for the 'remote' endpoint
type MsgRemoteResponse struct {
	Msg        string      `json:"msg"`
	FromRemote MsgResponse `json:"fromRemote"`
}

// ErrorResponse is an HTTP response in case of error
type ErrorResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}
