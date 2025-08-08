CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_collector_managements_status ON collector_managements(status);
CREATE INDEX IF NOT EXISTS idx_waste_drop_requests_status ON waste_drop_requests(status);
CREATE INDEX IF NOT EXISTS idx_waste_bank_priced_types ON waste_bank_priced_types(waste_bank_id, waste_type_id);
CREATE INDEX idx_users_email_change_token ON users(email_change_token);
CREATE INDEX idx_users_email_change_expiry ON users(email_change_expiry);