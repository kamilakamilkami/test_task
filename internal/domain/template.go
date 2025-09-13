package domain

import (
	"time"
	"github.com/google/uuid"
)

type Template struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`   
	Engine      string    `json:"engine"` 
	Body        string    `json:"body"`   
	Placeholders string   `json:"placeholders"` 
	Version     int       `json:"version"`
	IsActive    bool      `json:"isActive"`
	CreatedBy   uuid.UUID `json:"createdBy"`
	CreatedAt   time.Time `json:"createdAt"`
}

type PreviewRequest struct {
    Data map[string]interface{} `json:"data"`
}
