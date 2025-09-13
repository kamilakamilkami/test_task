ALTER TABLE users
ADD COLUMN employee_id UUID UNIQUE, 
ADD CONSTRAINT fk_users_employee
    FOREIGN KEY (employee_id) REFERENCES employees(id)
    ON DELETE CASCADE;
