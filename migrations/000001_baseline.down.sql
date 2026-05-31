DROP TRIGGER IF EXISTS audit_log_no_delete ON audit.audit_log;
DROP TRIGGER IF EXISTS audit_log_no_update ON audit.audit_log;
DROP FUNCTION IF EXISTS audit.deny_audit_log_mutation();

DROP INDEX IF EXISTS audit.idx_audit_log_resource;
DROP INDEX IF EXISTS audit.idx_audit_log_occurred_at;
DROP TABLE IF EXISTS audit.audit_log;

DROP SCHEMA IF EXISTS audit;
DROP SCHEMA IF EXISTS core;
