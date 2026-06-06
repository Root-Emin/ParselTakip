-- +goose Up
-- Sistem rolleri: Belediye Yöneticisi, Belediye Personeli, Müteahhit, Hak Sahibi, Sistem Yöneticisi
CREATE TABLE IF NOT EXISTS system_role_definitions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code            VARCHAR(50) NOT NULL UNIQUE,
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    permissions     JSONB NOT NULL DEFAULT '[]',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_system_role_definitions_code ON system_role_definitions(code);

-- +goose Down
DROP TABLE IF EXISTS system_role_definitions;
