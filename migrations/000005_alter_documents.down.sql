DROP TABLE IF EXISTS documents;

CREATE TABLE documents (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    version INT NOT NULL DEFAULT 1
);
