package bank

type AccountResponse struct {
	Status  string  `json:"status"`
	Account Account `json:"account"`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}
