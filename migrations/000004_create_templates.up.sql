CREATE TABLE templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,                        
    engine VARCHAR(10) NOT NULL,
    body TEXT NOT NULL,                               
    placeholders JSONB NOT NULL,                      
    version INT NOT NULL DEFAULT 1,                   
    is_active BOOLEAN NOT NULL DEFAULT true,         
    created_by UUID NOT NULL,                        
    created_at TIMESTAMP NOT NULL DEFAULT now(),

    CONSTRAINT fk_created_by FOREIGN KEY (created_by) REFERENCES employees(id) ON DELETE SET NULL
);
