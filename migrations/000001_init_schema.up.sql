CREATE TABLE departments (
                             id UUID PRIMARY KEY,
                             name TEXT NOT NULL
);

CREATE TABLE positions (
                           id UUID PRIMARY KEY,
                           title TEXT NOT NULL,
                           salary NUMERIC(12,2) NOT NULL
);

CREATE TABLE employees (
                           id UUID PRIMARY KEY,
                           first_name TEXT NOT NULL,
                           last_name TEXT NOT NULL,
                           middle_name TEXT,
                           department_id UUID REFERENCES departments(id) ON DELETE SET NULL,
                           position_id UUID REFERENCES positions(id) ON DELETE SET NULL
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
