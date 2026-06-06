-- +goose Up
CREATE TABLE IF NOT EXISTS document_types (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code                    VARCHAR(100) NOT NULL UNIQUE,
    name                    VARCHAR(255) NOT NULL,
    description             TEXT,
    category                VARCHAR(50) NOT NULL CHECK (category IN ('legal', 'identity', 'contract', 'permit', 'tax', 'financial', 'technical', 'other')),
    is_mandatory            BOOLEAN NOT NULL DEFAULT TRUE,
    requires_notary         BOOLEAN NOT NULL DEFAULT FALSE,
    requires_municipal_stamp BOOLEAN NOT NULL DEFAULT FALSE,
    is_valid_without_notary BOOLEAN NOT NULL DEFAULT TRUE,
    invalid_reason          TEXT,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_document_types_code ON document_types(code);
CREATE INDEX idx_document_types_category ON document_types(category);

-- +goose Down
DROP TABLE IF EXISTS document_types;
