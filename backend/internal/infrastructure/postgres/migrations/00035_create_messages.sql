-- +goose Up
CREATE TABLE IF NOT EXISTS messages (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id     UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    app_id              UUID NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    project_id          UUID REFERENCES urban_transformation_projects(id) ON DELETE SET NULL,
    sender_id           UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    recipient_id        UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    parent_id           UUID REFERENCES messages(id) ON DELETE SET NULL,
    subject             VARCHAR(255),
    body                TEXT NOT NULL,
    is_read             BOOLEAN NOT NULL DEFAULT FALSE,
    read_at             TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_messages_recipient ON messages(recipient_id, is_read);
CREATE INDEX idx_messages_sender ON messages(sender_id);
CREATE INDEX idx_messages_project ON messages(project_id);
CREATE INDEX idx_messages_org_app ON messages(organization_id, app_id);

-- +goose Down
DROP TABLE IF EXISTS messages;
