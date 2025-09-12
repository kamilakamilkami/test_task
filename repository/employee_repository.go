package repository

import (
	"context"
	"project/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EmployeeRepository interface {
	Create(ctx context.Context, emp *domain.Employee) error
	GetAll(ctx context.Context) ([]domain.Employee, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Employee, error)
	Update(ctx context.Context, emp *domain.Employee) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type employeeRepository struct {
	db *pgxpool.Pool
}

func NewEmployeeRepository(db *pgxpool.Pool) EmployeeRepository {
	return &employeeRepository{db: db}
}

func (r *employeeRepository) Create(ctx context.Context, emp *domain.Employee) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO employees 
		(id, fio, iin, email, phone, birth_date, employed_at, terminated_at, status,
		 department_id, position_id, grade, employment_type, salary_base, salary_currency,
		 work_schedule, manager_id) 
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)`,
		emp.ID, emp.Fio, emp.IIN, emp.Email, emp.Phone, emp.BirthDate, emp.EmployedAt, emp.TerminatedAt,
		emp.Status, emp.DepartmentID, emp.PositionID, emp.Grade, emp.EmploymentType,
		emp.SalaryBase, emp.SalaryCurrency, emp.WorkSchedule, emp.ManagerID,
	)
	return err
}

func (r *employeeRepository) GetAll(ctx context.Context) ([]domain.Employee, error) {
	rows, err := r.db.Query(ctx, `SELECT id, fio, iin, email, phone, birth_date, employed_at, terminated_at,
		status, department_id, position_id, grade, employment_type, salary_base, salary_currency, 
		work_schedule, manager_id FROM employees`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []domain.Employee
	for rows.Next() {
		var emp domain.Employee
		err := rows.Scan(
			&emp.ID, &emp.Fio, &emp.IIN, &emp.Email, &emp.Phone, &emp.BirthDate, &emp.EmployedAt,
			&emp.TerminatedAt, &emp.Status, &emp.DepartmentID, &emp.PositionID, &emp.Grade,
			&emp.EmploymentType, &emp.SalaryBase, &emp.SalaryCurrency, &emp.WorkSchedule, &emp.ManagerID,
		)
		if err != nil {
			return nil, err
		}
		employees = append(employees, emp)
	}
	return employees, nil
}

func (r *employeeRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Employee, error) {
	var emp domain.Employee
	err := r.db.QueryRow(ctx, `SELECT id, fio, iin, email, phone, birth_date, employed_at, terminated_at,
		status, department_id, position_id, grade, employment_type, salary_base, salary_currency, 
		work_schedule, manager_id FROM employees WHERE id=$1`, id).Scan(
		&emp.ID, &emp.Fio, &emp.IIN, &emp.Email, &emp.Phone, &emp.BirthDate, &emp.EmployedAt,
		&emp.TerminatedAt, &emp.Status, &emp.DepartmentID, &emp.PositionID, &emp.Grade,
		&emp.EmploymentType, &emp.SalaryBase, &emp.SalaryCurrency, &emp.WorkSchedule, &emp.ManagerID,
	)
	if err != nil {
		return nil, err
	}
	return &emp, nil
}

func (r *employeeRepository) Update(ctx context.Context, emp *domain.Employee) error {
	_, err := r.db.Exec(ctx, `UPDATE employees SET fio=$1, iin=$2, email=$3, phone=$4, birth_date=$5, 
		employed_at=$6, terminated_at=$7, status=$8, department_id=$9, position_id=$10, grade=$11, 
		employment_type=$12, salary_base=$13, salary_currency=$14, work_schedule=$15, manager_id=$16 
		WHERE id=$17`,
		emp.Fio, emp.IIN, emp.Email, emp.Phone, emp.BirthDate, emp.EmployedAt, emp.TerminatedAt,
		emp.Status, emp.DepartmentID, emp.PositionID, emp.Grade, emp.EmploymentType,
		emp.SalaryBase, emp.SalaryCurrency, emp.WorkSchedule, emp.ManagerID, emp.ID,
	)
	return err
}

func (r *employeeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM employees WHERE id=$1`, id)
	return err
}
