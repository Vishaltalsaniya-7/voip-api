package request

type CallRequest struct {
	Caller string `json:"caller" binding:"required"`
	Callee string `json:"callee" binding:"required"`
}
