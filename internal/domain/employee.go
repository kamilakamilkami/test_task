package domain

import (
	"github.com/google/uuid"
	"time"
)

type EmployeeStatus string
type EmploymentType string

const (
	StatusActive     EmployeeStatus = "ACTIVE"
	StatusOnLeave    EmployeeStatus = "ON_LEAVE"
	StatusTerminated EmployeeStatus = "TERMINATED"

	TypeFull EmploymentType = "FULL"
	TypePart EmploymentType = "PART"
)

type Employee struct {
	ID             uuid.UUID      `json:"id"`
	Fio            string         `json:"fullName"`
	IIN            string         `json:"iin"`
	Email          string         `json:"email"`
	Phone          string         `json:"phone"`
	BirthDate      time.Time      `json:"birthDate"`
	EmployedAt     time.Time      `json:"employedAt"`
	TerminatedAt   *time.Time     `json:"terminatedAt,omitempty"`
	Status         EmployeeStatus `json:"status"`
	DepartmentID   uuid.UUID      `json:"departmentId"`
	PositionID     uuid.UUID      `json:"positionId"`
	Grade          string         `json:"grade"`
	EmploymentType EmploymentType `json:"employmentType"`
	SalaryBase     float64        `json:"salaryBase"`
	SalaryCurrency string         `json:"salaryCurrency"`
	WorkSchedule   string         `json:"workSchedule"`
	ManagerID      *uuid.UUID     `json:"managerId,omitempty"`
}
