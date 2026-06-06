-- +goose Up
CREATE TABLE IF NOT EXISTS rent_assistance_applications (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id         UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    app_id                  UUID NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    project_id              UUID NOT NULL REFERENCES urban_transformation_projects(id) ON DELETE CASCADE,
    owner_id                UUID NOT NULL REFERENCES property_owners(id) ON DELETE CASCADE,
    application_date        DATE NOT NULL DEFAULT CURRENT_DATE,
    ministry_reference_no   VARCHAR(100),
    iban                    VARCHAR(34) NOT NULL,
    monthly_rent_amount     DECIMAL(12, 2),
    lease_contract_ref      VARCHAR(255),
    lease_start_date        DATE,
    lease_end_date          DATE,
    document_id             UUID REFERENCES documents(id) ON DELETE SET NULL,
    status                  VARCHAR(50) NOT NULL DEFAULT 'submitted' CHECK (status IN (
        'submitted', 'under_review', 'approved', 'rejected', 'payment_started', 'completed'
    )),
    approved_amount         DECIMAL(12, 2),
    payment_status          VARCHAR(50) DEFAULT 'pending' CHECK (payment_status IN ('pending', 'partial', 'completed', 'stopped')),
    reviewed_by             UUID REFERENCES users(id) ON DELETE SET NULL,
    reviewed_at             TIMESTAMPTZ,
    rejection_reason        TEXT,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_rent_assistance_project ON rent_assistance_applications(project_id);
CREATE INDEX idx_rent_assistance_owner ON rent_assistance_applications(owner_id);
CREATE INDEX idx_rent_assistance_status ON rent_assistance_applications(status);

-- +goose Down
DROP TABLE IF EXISTS rent_assistance_applications;
