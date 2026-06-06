-- +goose Up
CREATE TABLE IF NOT EXISTS workflow_document_requirements (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_step_id    UUID NOT NULL REFERENCES workflow_step_definitions(id) ON DELETE CASCADE,
    document_type_id    UUID NOT NULL REFERENCES document_types(id) ON DELETE CASCADE,
    is_mandatory        BOOLEAN NOT NULL DEFAULT TRUE,
    responsible_role    VARCHAR(50) NOT NULL CHECK (responsible_role IN (
        'municipality_admin', 'municipality_staff', 'contractor', 'property_owner', 'system_admin'
    )),
    notes               TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(workflow_step_id, document_type_id)
);

CREATE INDEX idx_wf_doc_req_step ON workflow_document_requirements(workflow_step_id);
CREATE INDEX idx_wf_doc_req_doc ON workflow_document_requirements(document_type_id);

-- +goose Down
DROP TABLE IF EXISTS workflow_document_requirements;
