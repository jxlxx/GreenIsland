package bank

import "github.com/google/uuid"

type NewAccountPayload struct {
	UserID uuid.UUID `json:"user_id"`
}
