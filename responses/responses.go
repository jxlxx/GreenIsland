package responses

type Status string

const (
	OK    Status = "OK"
	Error Status = "400 - Bad Request"
)

type NewBankAccount struct {
	Status Status `json:"status"`
}

type Response struct {
	Status Status `json:"status"`
}
