CREATE TABLE IF NOT EXISTS rooms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT UNIQUE,
    capacity INTEGER,
    description TEXT
);

CREATE INDEX IF NOT EXISTS idx_rooms_name ON rooms(name);