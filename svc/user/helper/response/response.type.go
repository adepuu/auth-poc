package response

type BaseResponse struct {
	Timestamp int64  `json:"timestamp,omitempty"`
	Status    string `json:"status"`
	Message   string `json:"message,omitempty"`
	Data      any    `json:"data,omitempty"`
}
