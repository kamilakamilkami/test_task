package domain

import (
	"time"
	"github.com/google/uuid"
)

type Document struct {
	ID                uuid.UUID  `json:"id"`
	Type              string     `json:"type"`
	EmployeeID        uuid.UUID  `json:"employeeId"`
	TemplateID			uuid.UUID        `json:"templateId"`
	TemplateVersionID int        `json:"templateVersionId"`
	Number            string     `json:"number"`
	Date              time.Time  `json:"date"`
	Status            string     `json:"status"`
	FileID            *uuid.UUID `json:"fileId,omitempty"`
	Data              string     `json:"data"`
	Meta             string     `json:"meta"`
}
