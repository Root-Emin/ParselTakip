-- +goose Up
CREATE TABLE IF NOT EXISTS risk_assessment_reports (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id         UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    app_id                  UUID NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    project_id              UUID NOT NULL REFERENCES urban_transformation_projects(id) ON DELETE CASCADE,
    building_id             UUID NOT NULL REFERENCES buildings(id) ON DELETE CASCADE,
    report_no               VARCHAR(100) NOT NULL,
    issued_by_institution   VARCHAR(255) NOT NULL,
    issued_at               DATE NOT NULL,
    risk_level              VARCHAR(50) NOT NULL CHECK (risk_level IN ('low', 'medium', 'high', 'very_high')),
    is_risky                BOOLEAN NOT NULL DEFAULT TRUE,
    registry_notified_at    TIMESTAMPTZ,
    municipality_notified_at TIMESTAMPTZ,
    document_id             UUID REFERENCES documents(id) ON DELETE SET NULL,
    status                  VARCHAR(50) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'objected', 'confirmed', 'cancelled')),
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(project_id, building_id, report_no)
);

CREATE INDEX idx_risk_reports_project ON risk_assessment_reports(project_id);
CREATE INDEX idx_risk_reports_building ON risk_assessment_reports(building_id);
CREATE INDEX idx_risk_reports_status ON risk_assessment_reports(status);

CREATE TABLE IF NOT EXISTS land_registry_notifications (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id         UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    app_id                  UUID NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    project_id              UUID NOT NULL REFERENCES urban_transformation_projects(id) ON DELETE CASCADE,
    building_id             UUID NOT NULL REFERENCES buildings(id) ON DELETE CASCADE,
    risk_report_id          UUID REFERENCES risk_assessment_reports(id) ON DELETE SET NULL,
    registry_office         VARCHAR(255) NOT NULL,
    reference_no            VARCHAR(100),
    notified_at             TIMESTAMPTZ NOT NULL,
    notification_deadline   TIMESTAMPTZ,
    status                  VARCHAR(50) NOT NULL DEFAULT 'submitted' CHECK (status IN ('submitted', 'registered', 'objection_period', 'confirmed')),
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_land_registry_project ON land_registry_notifications(project_id);
CREATE INDEX idx_land_registry_building ON land_registry_notifications(building_id);

CREATE TABLE IF NOT EXISTS owner_objections (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id     UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    app_id              UUID NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    project_id          UUID NOT NULL REFERENCES urban_transformation_projects(id) ON DELETE CASCADE,
    owner_id            UUID NOT NULL REFERENCES property_owners(id) ON DELETE CASCADE,
    land_registry_id    UUID REFERENCES land_registry_notifications(id) ON DELETE SET NULL,
    objection_text      TEXT NOT NULL,
    filed_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    objection_deadline  TIMESTAMPTZ NOT NULL,
    status              VARCHAR(50) NOT NULL DEFAULT 'filed' CHECK (status IN ('filed', 'under_review', 'accepted', 'rejected')),
    resolution          TEXT,
    resolved_at         TIMESTAMPTZ,
    resolved_by         UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_owner_objections_project ON owner_objections(project_id);
CREATE INDEX idx_owner_objections_owner ON owner_objections(owner_id);
CREATE INDEX idx_owner_objections_status ON owner_objections(status);

-- +goose Down
DROP TABLE IF EXISTS owner_objections;
DROP TABLE IF EXISTS land_registry_notifications;
DROP TABLE IF EXISTS risk_assessment_reports;
