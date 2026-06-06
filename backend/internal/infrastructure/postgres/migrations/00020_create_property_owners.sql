-- +goose Up
CREATE TABLE IF NOT EXISTS property_owners (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id     UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    app_id              UUID NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    user_id             UUID REFERENCES users(id) ON DELETE SET NULL,
    unit_id             UUID NOT NULL REFERENCES building_units(id) ON DELETE CASCADE,
    first_name          VARCHAR(255) NOT NULL,
    last_name           VARCHAR(255) NOT NULL,
    identity_number     VARCHAR(11),
    phone               VARCHAR(30),
    email               VARCHAR(255),
    address             TEXT,
    iban                VARCHAR(34),
    ownership_ratio     DECIMAL(5, 4) NOT NULL DEFAULT 1.0000 CHECK (ownership_ratio > 0 AND ownership_ratio <= 1),
    is_primary_contact  BOOLEAN NOT NULL DEFAULT FALSE,
    status              VARCHAR(50) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'objection_filed', 'consent_given', 'archived')),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_property_owners_unit ON property_owners(unit_id);
CREATE INDEX idx_property_owners_user ON property_owners(user_id);
CREATE INDEX idx_property_owners_org ON property_owners(organization_id);
CREATE INDEX idx_property_owners_app ON property_owners(app_id);
CREATE INDEX idx_property_owners_identity ON property_owners(identity_number);

-- +goose Down
DROP TABLE IF EXISTS property_owners;
