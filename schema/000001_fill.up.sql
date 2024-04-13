CREATE TABLE IF NOT EXISTS "user" (
    id UUID NOT NULL PRIMARY KEY,
    username TEXT CONSTRAINT  username_length CHECK ( char_length(username) <= 30 ) NOT NULL,
    password_hash TEXT CONSTRAINT passwordHash_length CHECK ( char_length(password_hash) <= 64) NOT NULL,
    is_admin BOOLEAN NOT NULL DEFAULT false
    );