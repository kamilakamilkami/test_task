ALTER TABLE documents
ADD COLUMN department_id UUID,
ADD CONSTRAINT fk_department
    FOREIGN KEY (department_id) REFERENCES departments(id)
    ON DELETE RESTRICT;
