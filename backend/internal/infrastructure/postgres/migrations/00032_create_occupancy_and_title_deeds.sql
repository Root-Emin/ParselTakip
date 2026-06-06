-- +goose Up
CREATE TABLE IF NOT EXISTS occupancy_permits (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id     UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    app_id              UUID NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    project_id          UUID NOT NULL REFERENCES urban_transformation_projects(id) ON DELETE CASCADE,
    building_id         UUID NOT NULL REFERENCES buildings(id) ON DELETE CASCADE,
    permit_no           VARCHAR(100) NOT NULL,
    issued_at           DATE NOT NULL,
    document_id         UUID REFERENCES documents(id) ON DELETE SET NULL,
    approved_by         UUID REFERENCES users(id) ON DELETE SET NULL,
    status              VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'issued', 'rejected')),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_occupancy_permits_project ON occupancy_permits(project_id);
CREATE INDEX idx_occupancy_permits_building ON occupancy_permits(building_id);

CREATE TABLE IF NOT EXISTS title_deed_transfers (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id     UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    app_id              UUID NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    project_id          UUID NOT NULL REFERENCES urban_transformation_projects(id) ON DELETE CASCADE,
    unit_id             UUID NOT NULL REFERENCES building_units(id) ON DELETE CASCADE,
    owner_id            UUID NOT NULL REFERENCES property_owners(id) ON DELETE CASCADE,
    old_title_deed_no   VARCHAR(100),
    new_title_deed_no   VARCHAR(100) NOT NULL,
    transferred_at      DATE NOT NULL,
    document_id         UUID REFERENCES documents(id) ON DELETE SET NULL,
    status              VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'completed', 'disputed')),
    delivered_at        TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_title_deed_transfers_project ON title_deed_transfers(project_id);
CREATE INDEX idx_title_deed_transfers_unit ON title_deed_transfers(unit_id);
CREATE INDEX idx_title_deed_transfers_owner ON title_deed_transfers(owner_id);

-- +goose Down
DROP TABLE IF EXISTS title_deed_transfers;
DROP TABLE IF EXISTS occupancy_permits;
