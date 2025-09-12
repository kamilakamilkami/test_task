package domain

import "github.com/google/uuid"

type Department struct {
	ID       uuid.UUID  `json:"id"`
	Name     string     `json:"name"`
	Code     string     `json:"code"`
	ParentID *uuid.UUID `json:"parentId,omitempty"`
}
