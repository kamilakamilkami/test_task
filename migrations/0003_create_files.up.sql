CREATE TABLE files (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    size BIGINT NOT NULL,
    path TEXT NOT NULL,
    base64 TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
