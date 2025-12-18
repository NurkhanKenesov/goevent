-- +migrate Up
DROP TYPE IF EXISTS invitation_status;
CREATE TYPE invitation_status AS ENUM ('Pending', 'Accepted', 'Declined');

CREATE TABLE invitations (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    event_id INTEGER NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    status invitation_status DEFAULT 'Pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Performance indexes: prevent N+1 queries, speed up lookups
CREATE INDEX idx_invitations_user_id ON invitations(user_id);
CREATE INDEX idx_invitations_event_id ON invitations(event_id);
CREATE INDEX idx_invitations_user_status ON invitations(user_id, status);
CREATE INDEX idx_invitations_event_status ON invitations(event_id, status);

-- +migrate Down
DROP TABLE invitations;
DROP TYPE invitation_status;
