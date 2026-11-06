CREATE TABLE IF NOT EXISTS tickets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    category TEXT NOT NULL,
    priority TEXT NOT NULL DEFAULT 'medium' CHECK(priority IN('low', 'medium', 'high', 'urgent')),
    status TEXT NOT NULL DEFAULT 'open' CHECK(status IN('open', 'in_progress', 'resolved', 'closed')),
    requester_id UUID NOT NULL REFERENCES users(id),
    assignee_id UUID REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_tickets_status ON tickets(status);
CREATE INDEX IF NOT EXISTS idx_tickets_priority ON tickets(priority);
CREATE INDEX IF NOT EXISTS idx_tickets_requester_id ON tickets(requester_id);
CREATE INDEX IF NOT EXISTS idx_tickets_assignee_id ON tickets(assignee_id);
CREATE INDEX IF NOT EXISTS idx_tickets_created_at ON tickets(created_at DESC);
