package domain

import "github.com/google/uuid"

type Position struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Code string    `json:"code"`
}
