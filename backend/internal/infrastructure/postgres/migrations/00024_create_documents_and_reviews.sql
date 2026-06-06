-- +goose Up
CREATE TABLE IF NOT EXISTS documents (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id     UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    app_id              UUID NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    project_id          UUID NOT NULL REFERENCES urban_transformation_projects(id) ON DELETE CASCADE,
    document_type_id    UUID NOT NULL REFERENCES document_types(id) ON DELETE RESTRICT,
    building_id         UUID REFERENCES buildings(id) ON DELETE SET NULL,
    unit_id             UUID REFERENCES building_units(id) ON DELETE SET NULL,
    owner_id            UUID REFERENCES property_owners(id) ON DELETE SET NULL,
    file_name           VARCHAR(500) NOT NULL,
    file_path           TEXT NOT NULL,
    file_size           BIGINT,
    mime_type           VARCHAR(100),
    status              VARCHAR(50) NOT NULL DEFAULT 'draft' CHECK (status IN (
        'draft', 'submitted', 'under_review', 'approved', 'rejected', 'missing', 'expired', 'invalid'
    )),
    is_notarized        BOOLEAN NOT NULL DEFAULT FALSE,
    notary_date         DATE,
    expiry_date         DATE,
    uploaded_by         UUID REFERENCES users(id) ON DELETE SET NULL,
    uploaded_by_role    VARCHAR(50),
    version             INT NOT NULL DEFAULT 1,
    metadata            JSONB,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_documents_project ON documents(project_id);
CREATE INDEX idx_documents_type ON documents(document_type_id);
CREATE INDEX idx_documents_status ON documents(status);
CREATE INDEX idx_documents_owner ON documents(owner_id);
CREATE INDEX idx_documents_building ON documents(building_id);
CREATE INDEX idx_documents_org_app ON documents(organization_id, app_id);

CREATE TABLE IF NOT EXISTS document_reviews (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id     UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    app_id              UUID NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    document_id         UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    reviewer_id         UUID NOT NULL REFERENCES users(id) ON DELETE SET NULL,
    status              VARCHAR(50) NOT NULL CHECK (status IN ('approved', 'rejected', 'missing_items', 'needs_revision')),
    missing_items       JSONB,
    review_notes        TEXT,
    reviewed_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_document_reviews_document ON document_reviews(document_id);
CREATE INDEX idx_document_reviews_reviewer ON document_reviews(reviewer_id);
CREATE INDEX idx_document_reviews_status ON document_reviews(status);

-- +goose Down
DROP TABLE IF EXISTS document_reviews;
DROP TABLE IF EXISTS documents;
