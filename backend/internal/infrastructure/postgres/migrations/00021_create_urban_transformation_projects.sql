-- +goose Up
CREATE TABLE IF NOT EXISTS urban_transformation_projects (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id         UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    app_id                  UUID NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    code                    VARCHAR(50) NOT NULL,
    name                    VARCHAR(255) NOT NULL,
    description             TEXT,
    status                  VARCHAR(50) NOT NULL DEFAULT 'draft' CHECK (status IN (
        'draft', 'initiated', 'in_progress', 'blocked', 'completed', 'cancelled'
    )),
    current_workflow_step_id UUID REFERENCES workflow_step_definitions(id) ON DELETE SET NULL,
    initiated_by            UUID REFERENCES users(id) ON DELETE SET NULL,
    assigned_contractor_id  UUID REFERENCES contractor_companies(id) ON DELETE SET NULL,
    started_at              TIMESTAMPTZ,
    target_completion_at    TIMESTAMPTZ,
    completed_at            TIMESTAMPTZ,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(organization_id, code)
);

CREATE INDEX idx_ut_projects_org ON urban_transformation_projects(organization_id);
CREATE INDEX idx_ut_projects_app ON urban_transformation_projects(app_id);
CREATE INDEX idx_ut_projects_status ON urban_transformation_projects(status);
CREATE INDEX idx_ut_projects_step ON urban_transformation_projects(current_workflow_step_id);
CREATE INDEX idx_ut_projects_contractor ON urban_transformation_projects(assigned_contractor_id);

-- +goose Down
DROP TABLE IF EXISTS urban_transformation_projects;
