WITH cte AS (
    SELECT ticket_id
    FROM tickets
    WHERE event_id = $1 AND status = 'AVAILABLE'
    ORDER BY ticket_id
    FOR UPDATE SKIP LOCKED
    LIMIT 1
)
UPDATE tickets
SET status = 'HELD', held_by = $2, hold_expiry = now() + make_interval(mins => $3)
FROM cte
WHERE tickets.ticket_id = cte.ticket_id
RETURNING tickets.ticket_id;
