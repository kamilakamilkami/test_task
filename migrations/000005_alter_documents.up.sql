DROP TABLE IF EXISTS documents;

CREATE TABLE documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type VARCHAR(30) NOT NULL,                
    employee_id UUID NOT NULL,                
    template_version_id INT NOT NULL,         
    number VARCHAR(50) NOT NULL UNIQUE,       
    date TIMESTAMP NOT NULL,                  
    status VARCHAR(30) NOT NULL,              
    file_id UUID,                             
    data JSONB NOT NULL,                      
    meta JSONB,                               

    CONSTRAINT fk_employee FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE,
    CONSTRAINT fk_template_version FOREIGN KEY (template_version_id) REFERENCES template_versions(id) ON DELETE RESTRICT,
    CONSTRAINT fk_file FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE SET NULL
);
