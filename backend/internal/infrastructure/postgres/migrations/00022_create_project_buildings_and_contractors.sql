-- +goose Up
CREATE TABLE IF NOT EXISTS project_buildings (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id     UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    app_id              UUID NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    project_id          UUID NOT NULL REFERENCES urban_transformation_projects(id) ON DELETE CASCADE,
    building_id         UUID NOT NULL REFERENCES buildings(id) ON DELETE CASCADE,
    added_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    added_by            UUID REFERENCES users(id) ON DELETE SET NULL,
    UNIQUE(project_id, building_id)
);

CREATE INDEX idx_project_buildings_project ON project_buildings(project_id);
CREATE INDEX idx_project_buildings_building ON project_buildings(building_id);
CREATE INDEX idx_project_buildings_org ON project_buildings(organization_id);

CREATE TABLE IF NOT EXISTS project_contractor_assignments (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id         UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    app_id                  UUID NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    project_id              UUID NOT NULL REFERENCES urban_transformation_projects(id) ON DELETE CASCADE,
    contractor_company_id   UUID NOT NULL REFERENCES contractor_companies(id) ON DELETE CASCADE,
    assigned_by             UUID NOT NULL REFERENCES users(id) ON DELETE SET NULL,
    assigned_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    status                  VARCHAR(50) NOT NULL DEFAULT 'active' CHECK (status IN ('pending', 'active', 'completed', 'terminated')),
    contract_share_notes    TEXT,
    terminated_at           TIMESTAMPTZ,
    termination_reason      TEXT,
    UNIQUE(project_id, contractor_company_id)
);

CREATE INDEX idx_project_contractors_project ON project_contractor_assignments(project_id);
CREATE INDEX idx_project_contractors_company ON project_contractor_assignments(contractor_company_id);
CREATE INDEX idx_project_contractors_status ON project_contractor_assignments(status);

-- +goose Down
DROP TABLE IF EXISTS project_contractor_assignments;
DROP TABLE IF EXISTS project_buildings;
