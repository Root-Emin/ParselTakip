-- +goose Up
CREATE TABLE IF NOT EXISTS building_units (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id     UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    app_id              UUID NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    building_id         UUID NOT NULL REFERENCES buildings(id) ON DELETE CASCADE,
    unit_no             VARCHAR(50) NOT NULL,
    floor_no            INT,
    area_sqm            DECIMAL(10, 2),
    room_count          VARCHAR(20),
    ownership_type      VARCHAR(50) NOT NULL DEFAULT 'kat_mulkiyeti' CHECK (ownership_type IN ('kat_irtifaki', 'kat_mulkiyeti', 'arsa_payi')),
    title_deed_no       VARCHAR(100),
    status              VARCHAR(50) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'in_transformation', 'transferred', 'archived')),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(building_id, unit_no)
);

CREATE INDEX idx_building_units_building ON building_units(building_id);
CREATE INDEX idx_building_units_org ON building_units(organization_id);
CREATE INDEX idx_building_units_app ON building_units(app_id);
CREATE INDEX idx_building_units_title_deed ON building_units(title_deed_no);

-- +goose Down
DROP TABLE IF EXISTS building_units;
