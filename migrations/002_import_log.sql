-- Create migration log table
CREATE TABLE IF NOT EXISTS import_log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    started_at TIMESTAMP,
    ended_at TIMESTAMP,
    import_count INTEGER,
    total_remote_records INTEGER,
    success INTEGER,
    error_message TEXT
);
