CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE SCHEMA IF NOT EXISTS core;
CREATE SCHEMA IF NOT EXISTS audit;

CREATE TABLE audit.audit_log (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    occurred_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    actor_id      TEXT,
    action        TEXT NOT NULL,
    resource_type TEXT NOT NULL,
    resource_id   TEXT,
    metadata      JSONB NOT NULL DEFAULT '{}',
    request_id    UUID,
    ip_address    INET,
    user_agent    TEXT
);

CREATE INDEX idx_audit_log_occurred_at ON audit.audit_log (occurred_at DESC);
CREATE INDEX idx_audit_log_resource ON audit.audit_log (resource_type, resource_id);

CREATE OR REPLACE FUNCTION audit.deny_audit_log_mutation()
RETURNS TRIGGER AS $$
BEGIN
    RAISE EXCEPTION 'audit.audit_log is append-only: % operations are forbidden', TG_OP;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER audit_log_no_update
    BEFORE UPDATE ON audit.audit_log
    FOR EACH ROW EXECUTE FUNCTION audit.deny_audit_log_mutation();

CREATE TRIGGER audit_log_no_delete
    BEFORE DELETE ON audit.audit_log
    FOR EACH ROW EXECUTE FUNCTION audit.deny_audit_log_mutation();
