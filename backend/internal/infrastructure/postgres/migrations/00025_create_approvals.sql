-- +goose Up
CREATE TABLE IF NOT EXISTS approvals (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id     UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    app_id              UUID NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    project_id          UUID NOT NULL REFERENCES urban_transformation_projects(id) ON DELETE CASCADE,
    approval_type       VARCHAR(50) NOT NULL CHECK (approval_type IN (
        'municipal_initiation', 'owner_consent', 'majority_decision', 'municipal_permit',
        'demolition', 'occupancy', 'title_transfer', 'rent_assistance', 'contractor_assignment'
    )),
    approver_id         UUID REFERENCES users(id) ON DELETE SET NULL,
    approver_role       VARCHAR(50) NOT NULL,
    owner_id            UUID REFERENCES property_owners(id) ON DELETE SET NULL,
    status              VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected', 'expired')),
    decision_notes      TEXT,
    expires_at          TIMESTAMPTZ,
    decided_at          TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_approvals_project ON approvals(project_id);
CREATE INDEX idx_approvals_type ON approvals(approval_type);
CREATE INDEX idx_approvals_status ON approvals(status);
CREATE INDEX idx_approvals_approver ON approvals(approver_id);
CREATE INDEX idx_approvals_owner ON approvals(owner_id);

-- +goose Down
DROP TABLE IF EXISTS approvals;
