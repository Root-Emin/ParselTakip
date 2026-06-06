-- +goose Up
CREATE TABLE IF NOT EXISTS project_workflow_states (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id     UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    app_id              UUID NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    project_id          UUID NOT NULL REFERENCES urban_transformation_projects(id) ON DELETE CASCADE,
    workflow_step_id    UUID NOT NULL REFERENCES workflow_step_definitions(id) ON DELETE RESTRICT,
    status              VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (status IN (
        'pending', 'in_progress', 'awaiting_documents', 'awaiting_approval', 'completed', 'blocked', 'skipped'
    )),
    started_at          TIMESTAMPTZ,
    completed_at        TIMESTAMPTZ,
    due_at              TIMESTAMPTZ,
    blocked_reason      TEXT,
    updated_by          UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(project_id, workflow_step_id)
);

CREATE INDEX idx_project_wf_states_project ON project_workflow_states(project_id);
CREATE INDEX idx_project_wf_states_step ON project_workflow_states(workflow_step_id);
CREATE INDEX idx_project_wf_states_status ON project_workflow_states(status);

CREATE TABLE IF NOT EXISTS project_workflow_history (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id     UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    app_id              UUID NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    project_id          UUID NOT NULL REFERENCES urban_transformation_projects(id) ON DELETE CASCADE,
    from_step_id        UUID REFERENCES workflow_step_definitions(id) ON DELETE SET NULL,
    to_step_id          UUID NOT NULL REFERENCES workflow_step_definitions(id) ON DELETE RESTRICT,
    action              VARCHAR(100) NOT NULL,
    notes               TEXT,
    changed_by          UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_project_wf_history_project ON project_workflow_history(project_id);
CREATE INDEX idx_project_wf_history_created ON project_workflow_history(created_at);

-- +goose Down
DROP TABLE IF EXISTS project_workflow_history;
DROP TABLE IF EXISTS project_workflow_states;
