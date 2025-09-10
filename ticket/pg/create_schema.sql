CREATE TABLE IF NOT EXISTS tickets (
    ticket_id BIGSERIAL PRIMARY KEY,
    event_id BIGINT NOT NULL,
    status TEXT NOT NULL DEFAULT 'AVAILABLE',
    held_by TEXT,
    hold_expiry TIMESTAMPTZ
)
