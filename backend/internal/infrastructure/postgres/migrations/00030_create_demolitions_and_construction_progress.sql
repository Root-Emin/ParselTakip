-- +goose Up
CREATE TABLE IF NOT EXISTS demolitions (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id         UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    app_id                  UUID NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    project_id              UUID NOT NULL REFERENCES urban_transformation_projects(id) ON DELETE CASCADE,
    building_id             UUID NOT NULL REFERENCES buildings(id) ON DELETE CASCADE,
    contractor_company_id   UUID REFERENCES contractor_companies(id) ON DELETE SET NULL,
    permit_id               UUID REFERENCES building_permits(id) ON DELETE SET NULL,
    scheduled_date          DATE,
    started_at              TIMESTAMPTZ,
    completed_at            TIMESTAMPTZ,
    approved_by             UUID REFERENCES users(id) ON DELETE SET NULL,
    status                  VARCHAR(50) NOT NULL DEFAULT 'planned' CHECK (status IN ('planned', 'in_progress', 'completed', 'cancelled')),
    notes                   TEXT,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_demolitions_project ON demolitions(project_id);
CREATE INDEX idx_demolitions_building ON demolitions(building_id);
CREATE INDEX idx_demolitions_status ON demolitions(status);

CREATE TABLE IF NOT EXISTS construction_progress (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id     UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    app_id              UUID NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    project_id          UUID NOT NULL REFERENCES urban_transformation_projects(id) ON DELETE CASCADE,
    progress_percentage DECIMAL(5, 2) NOT NULL DEFAULT 0 CHECK (progress_percentage >= 0 AND progress_percentage <= 100),
    milestone           VARCHAR(255),
    notes               TEXT,
    reported_by         UUID REFERENCES users(id) ON DELETE SET NULL,
    reported_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_construction_progress_project ON construction_progress(project_id);
CREATE INDEX idx_construction_progress_reported ON construction_progress(reported_at);

-- +goose Down
DROP TABLE IF EXISTS construction_progress;
DROP TABLE IF EXISTS demolitions;
