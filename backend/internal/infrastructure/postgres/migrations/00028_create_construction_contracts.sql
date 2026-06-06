-- +goose Up
CREATE TABLE IF NOT EXISTS construction_contracts (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id         UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    app_id                  UUID NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    project_id              UUID NOT NULL REFERENCES urban_transformation_projects(id) ON DELETE CASCADE,
    contractor_company_id   UUID NOT NULL REFERENCES contractor_companies(id) ON DELETE RESTRICT,
    contract_no             VARCHAR(100),
    contract_type           VARCHAR(50) NOT NULL DEFAULT 'kat_karsiligi' CHECK (contract_type IN ('kat_karsiligi', 'cash', 'mixed')),
    owner_share_ratio       DECIMAL(5, 4),
    contractor_share_ratio  DECIMAL(5, 4),
    delivery_months         INT,
    delivery_deadline       DATE,
    is_notarized            BOOLEAN NOT NULL DEFAULT FALSE,
    notary_name             VARCHAR(255),
    notary_date             DATE,
    document_id             UUID REFERENCES documents(id) ON DELETE SET NULL,
    status                  VARCHAR(50) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'signed', 'active', 'completed', 'terminated', 'invalid')),
    invalid_reason          TEXT,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_construction_contracts_project ON construction_contracts(project_id);
CREATE INDEX idx_construction_contracts_contractor ON construction_contracts(contractor_company_id);
CREATE INDEX idx_construction_contracts_status ON construction_contracts(status);

-- +goose Down
DROP TABLE IF EXISTS construction_contracts;
