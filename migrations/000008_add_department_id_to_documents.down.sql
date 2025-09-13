ALTER TABLE documents
DROP CONSTRAINT IF EXISTS fk_department,
DROP COLUMN IF EXISTS department_id;
