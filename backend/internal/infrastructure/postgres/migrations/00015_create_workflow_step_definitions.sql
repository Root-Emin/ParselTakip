-- +goose Up
CREATE TABLE IF NOT EXISTS workflow_step_definitions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    step_order      INT NOT NULL UNIQUE,
    code            VARCHAR(100) NOT NULL UNIQUE,
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    responsible_role VARCHAR(50) NOT NULL CHECK (responsible_role IN (
        'municipality_admin', 'municipality_staff', 'contractor', 'property_owner', 'system_admin'
    )),
    sla_days        INT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_workflow_steps_order ON workflow_step_definitions(step_order);
CREATE INDEX idx_workflow_steps_code ON workflow_step_definitions(code);

-- +goose Down
DROP TABLE IF EXISTS workflow_step_definitions;
