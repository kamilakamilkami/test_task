CREATE TABLE document_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    workflow TEXT
);

INSERT INTO document_types (code, name, workflow)
VALUES (
    'CERT',
    'Справка с места работы',
    ''
);
