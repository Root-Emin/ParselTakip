-- +goose Up
CREATE TABLE IF NOT EXISTS contractor_companies (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id     UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    app_id              UUID NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    user_id             UUID REFERENCES users(id) ON DELETE SET NULL,
    company_name        VARCHAR(255) NOT NULL,
    tax_number          VARCHAR(20) NOT NULL,
    trade_registry_no   VARCHAR(50),
    authorized_person   VARCHAR(255),
    phone               VARCHAR(30),
    email               VARCHAR(255),
    address             TEXT,
    status              VARCHAR(50) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'suspended', 'blacklisted')),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(organization_id, tax_number)
);

CREATE INDEX idx_contractor_companies_org ON contractor_companies(organization_id);
CREATE INDEX idx_contractor_companies_app ON contractor_companies(app_id);
CREATE INDEX idx_contractor_companies_user ON contractor_companies(user_id);
CREATE INDEX idx_contractor_companies_status ON contractor_companies(status);

-- +goose Down
DROP TABLE IF EXISTS contractor_companies;
