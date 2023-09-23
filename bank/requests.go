package bank

import "github.com/google/uuid"

type NewAccountRequest struct {
	UserID uuid.UUID `json:"user_id"`
}
