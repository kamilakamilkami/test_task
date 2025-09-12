package domain

import "github.com/google/uuid"

type Doctype struct {
	ID       uuid.UUID `json:"id"`
	Code	 string 	`json:"code"`
	Name     string    `json:"name"`
	Workflow string		`json:"workflow"`
}
