CREATE TABLE IF NOT EXISTS number_sequences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    doc_type VARCHAR(50) NOT NULL UNIQUE,
    pattern VARCHAR(100) NOT NULL,
    last_value INT NOT NULL DEFAULT 0,
    updated_at TIMESTAMP DEFAULT now()
);
