-- +goose Up
CREATE TABLE IF NOT EXISTS building_permits (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id     UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    app_id              UUID NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    project_id          UUID NOT NULL REFERENCES urban_transformation_projects(id) ON DELETE CASCADE,
    building_id         UUID REFERENCES buildings(id) ON DELETE SET NULL,
    permit_type         VARCHAR(50) NOT NULL CHECK (permit_type IN (
        'imar_uygunluk', 'yapi_ruhsati', 'yikim_ruhsati', 'iskan'
    )),
    permit_no           VARCHAR(100) NOT NULL,
    issued_by           VARCHAR(255),
    issued_at           DATE NOT NULL,
    expires_at          DATE,
    document_id         UUID REFERENCES documents(id) ON DELETE SET NULL,
    approved_by           UUID REFERENCES users(id) ON DELETE SET NULL,
    status              VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected', 'expired', 'revoked')),
    rejection_reason    TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_building_permits_project ON building_permits(project_id);
CREATE INDEX idx_building_permits_type ON building_permits(permit_type);
CREATE INDEX idx_building_permits_status ON building_permits(status);

CREATE TABLE IF NOT EXISTS tax_exemption_documents (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id     UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    app_id              UUID NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    project_id          UUID NOT NULL REFERENCES urban_transformation_projects(id) ON DELETE CASCADE,
    exemption_type      VARCHAR(50) NOT NULL CHECK (exemption_type IN (
        'tapu_harci', 'damga_vergisi', 'noter_harci', 'other'
    )),
    reference_no        VARCHAR(100),
    issued_at           DATE,
    document_id         UUID REFERENCES documents(id) ON DELETE SET NULL,
    status              VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tax_exemptions_project ON tax_exemption_documents(project_id);

-- +goose Down
DROP TABLE IF EXISTS tax_exemption_documents;
DROP TABLE IF EXISTS building_permits;
