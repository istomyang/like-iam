package v1

type Response struct {
	Allowed bool   `json:"allowed"`
	Reason  string `json:"reason,omitempty"`
	Error   string `json:"error,omitempty"`
}
