-- +goose Up
CREATE TABLE IF NOT EXISTS appointments (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id     UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    app_id              UUID NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    project_id          UUID REFERENCES urban_transformation_projects(id) ON DELETE SET NULL,
    owner_id            UUID REFERENCES property_owners(id) ON DELETE SET NULL,
    title               VARCHAR(255) NOT NULL,
    description         TEXT,
    location            VARCHAR(500),
    scheduled_at        TIMESTAMPTZ NOT NULL,
    duration_minutes    INT,
    status              VARCHAR(50) NOT NULL DEFAULT 'scheduled' CHECK (status IN (
        'scheduled', 'completed', 'cancelled', 'rescheduled'
    )),
    created_by          UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_appointments_project ON appointments(project_id);
CREATE INDEX idx_appointments_owner ON appointments(owner_id);
CREATE INDEX idx_appointments_status ON appointments(status);
CREATE INDEX idx_appointments_scheduled ON appointments(scheduled_at);
CREATE INDEX idx_appointments_org_app ON appointments(organization_id, app_id);

-- +goose Down
DROP TABLE IF EXISTS appointments;
