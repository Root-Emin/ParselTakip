-- +goose Up
CREATE TABLE IF NOT EXISTS project_notifications (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id     UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    app_id              UUID NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    project_id          UUID REFERENCES urban_transformation_projects(id) ON DELETE CASCADE,
    user_id             UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    notification_type   VARCHAR(50) NOT NULL CHECK (notification_type IN (
        'document_required', 'document_approved', 'document_rejected', 'document_missing',
        'step_completed', 'step_blocked', 'approval_required', 'objection_deadline',
        'meeting_scheduled', 'contract_signed', 'permit_issued', 'demolition_scheduled',
        'rent_assistance_update', 'construction_update', 'title_deed_ready', 'general'
    )),
    title               VARCHAR(255) NOT NULL,
    message             TEXT NOT NULL,
    channel             VARCHAR(20) NOT NULL DEFAULT 'in_app' CHECK (channel IN ('in_app', 'email', 'sms')),
    is_read             BOOLEAN NOT NULL DEFAULT FALSE,
    read_at             TIMESTAMPTZ,
    metadata            JSONB,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_project_notifications_user ON project_notifications(user_id);
CREATE INDEX idx_project_notifications_project ON project_notifications(project_id);
CREATE INDEX idx_project_notifications_unread ON project_notifications(user_id, is_read) WHERE is_read = FALSE;
CREATE INDEX idx_project_notifications_type ON project_notifications(notification_type);

-- +goose Down
DROP TABLE IF EXISTS project_notifications;
