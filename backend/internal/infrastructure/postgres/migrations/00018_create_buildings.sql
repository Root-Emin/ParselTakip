-- +goose Up
CREATE TABLE IF NOT EXISTS buildings (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id     UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    app_id              UUID NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    name                VARCHAR(255),
    address             TEXT NOT NULL,
    city                VARCHAR(100) NOT NULL,
    district            VARCHAR(100) NOT NULL,
    neighborhood        VARCHAR(100),
    block_no            VARCHAR(50),
    parcel_no           VARCHAR(50),
    island_no           VARCHAR(50),
    floor_count         INT,
    unit_count          INT NOT NULL DEFAULT 1,
    construction_year   INT,
    building_type       VARCHAR(50) DEFAULT 'residential' CHECK (building_type IN ('residential', 'commercial', 'mixed')),
    risk_status         VARCHAR(50) NOT NULL DEFAULT 'unknown' CHECK (risk_status IN ('unknown', 'risky', 'not_risky', 'under_assessment')),
    status              VARCHAR(50) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'in_transformation', 'demolished', 'rebuilt', 'archived')),
    created_by          UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_buildings_org ON buildings(organization_id);
CREATE INDEX idx_buildings_app ON buildings(app_id);
CREATE INDEX idx_buildings_city_district ON buildings(city, district);
CREATE INDEX idx_buildings_risk_status ON buildings(risk_status);
CREATE INDEX idx_buildings_status ON buildings(status);

-- +goose Down
DROP TABLE IF EXISTS buildings;
