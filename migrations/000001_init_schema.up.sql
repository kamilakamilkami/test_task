CREATE TABLE departments (
                             id UUID PRIMARY KEY,
                             name VARCHAR(100) NOT NULL,
                             code VARCHAR(50) UNIQUE NOT NULL,
                             parent_id UUID NULL REFERENCES departments(id) ON DELETE SET NULL
);

CREATE TABLE positions (
                           id UUID PRIMARY KEY,
                           name VARCHAR(255) NOT NULL,
                           code VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE employees (
                           id UUID PRIMARY KEY,
                           fio VARCHAR(150) NOT NULL,
                           iin VARCHAR(12) UNIQUE NOT NULL,
                           email VARCHAR(150) UNIQUE NOT NULL,
                           phone VARCHAR(50),
                           birth_date DATE NOT NULL,
                           employed_at DATE NOT NULL,
                           terminated_at DATE,
                           status VARCHAR(20) NOT NULL CHECK (status IN ('ACTIVE','ON_LEAVE','TERMINATED')),
                           department_id UUID NOT NULL REFERENCES departments(id) ON DELETE CASCADE,
                           position_id UUID NOT NULL,
                           grade VARCHAR(50),
                           employment_type VARCHAR(10) NOT NULL CHECK (employment_type IN ('FULL','PART')),
                           salary_base NUMERIC(12,2),
                           salary_currency VARCHAR(10),
                           work_schedule VARCHAR(100),
                           manager_id UUID,
                           created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE contracts (
                           id UUID PRIMARY KEY,
                           employee_id UUID REFERENCES employees(id) ON DELETE CASCADE,
                           start_date DATE NOT NULL,
                           end_date DATE,
                           type TEXT NOT NULL
);

CREATE TABLE orders (
                        id UUID PRIMARY KEY,
                        number TEXT NOT NULL,
                        date DATE NOT NULL,
                        type TEXT NOT NULL,
                        employee_id UUID REFERENCES employees(id) ON DELETE CASCADE
);

CREATE TABLE documents (
                           id UUID PRIMARY KEY,
                           name TEXT NOT NULL,
                           content TEXT NOT NULL,
                           created_at TIMESTAMP NOT NULL DEFAULT now(),
                           version INT NOT NULL DEFAULT 1
);

CREATE TABLE users (
                       id UUID PRIMARY KEY,
                       name VARCHAR(100) NOT NULL,
                       email VARCHAR(100) UNIQUE NOT NULL,
                       password TEXT NOT NULL,
                       role VARCHAR(20) NOT NULL DEFAULT 'EMPLOYEE', -- ADMIN, HR, MANAGER, EMPLOYEE
                       created_at TIMESTAMP DEFAULT now(),
                       updated_at TIMESTAMP DEFAULT now()
);