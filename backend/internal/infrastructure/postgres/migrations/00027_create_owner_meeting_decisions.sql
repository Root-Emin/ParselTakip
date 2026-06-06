-- +goose Up
CREATE TABLE IF NOT EXISTS owner_meeting_decisions (
    id                          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id             UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    app_id                      UUID NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    project_id                  UUID NOT NULL REFERENCES urban_transformation_projects(id) ON DELETE CASCADE,
    meeting_date                TIMESTAMPTZ NOT NULL,
    total_ownership_ratio       DECIMAL(5, 4) NOT NULL,
    votes_for_ratio             DECIMAL(5, 4) NOT NULL DEFAULT 0,
    votes_against_ratio         DECIMAL(5, 4) NOT NULL DEFAULT 0,
    abstained_ratio             DECIMAL(5, 4) NOT NULL DEFAULT 0,
    quorum_required             DECIMAL(5, 4) NOT NULL DEFAULT 0.6667,
    quorum_met                  BOOLEAN NOT NULL DEFAULT FALSE,
    selected_contractor_id      UUID REFERENCES contractor_companies(id) ON DELETE SET NULL,
    decision_text               TEXT,
    document_id                 UUID REFERENCES documents(id) ON DELETE SET NULL,
    recorded_by                 UUID REFERENCES users(id) ON DELETE SET NULL,
    status                      VARCHAR(50) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'confirmed', 'invalid')),
    created_at                  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at                  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_owner_meeting_project ON owner_meeting_decisions(project_id);
CREATE INDEX idx_owner_meeting_quorum ON owner_meeting_decisions(quorum_met);

CREATE TABLE IF NOT EXISTS owner_meeting_votes (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    meeting_decision_id UUID NOT NULL REFERENCES owner_meeting_decisions(id) ON DELETE CASCADE,
    owner_id            UUID NOT NULL REFERENCES property_owners(id) ON DELETE CASCADE,
    vote                VARCHAR(20) NOT NULL CHECK (vote IN ('for', 'against', 'abstain')),
    ownership_ratio     DECIMAL(5, 4) NOT NULL,
    voted_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(meeting_decision_id, owner_id)
);

CREATE INDEX idx_owner_votes_meeting ON owner_meeting_votes(meeting_decision_id);
CREATE INDEX idx_owner_votes_owner ON owner_meeting_votes(owner_id);

-- +goose Down
DROP TABLE IF EXISTS owner_meeting_votes;
DROP TABLE IF EXISTS owner_meeting_decisions;
