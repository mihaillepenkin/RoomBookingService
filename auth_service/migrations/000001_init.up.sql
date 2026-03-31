CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT UNIQUE CHECK (char_length(email) <= 255),
    role TEXT NOT NULL,
    password TEXT NOT NULL
);