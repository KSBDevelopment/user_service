CREATE TABLE settings (
                          id SERIAL PRIMARY KEY,
                          user_id INTEGER NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
                          is_private BOOLEAN DEFAULT FALSE,
                          dark_mode BOOLEAN DEFAULT FALSE,
                          language VARCHAR(10) DEFAULT 'en',
                          created_at TIMESTAMPTZ DEFAULT NOW(),
                          updated_at TIMESTAMPTZ DEFAULT NOW(),
                          deleted_at TIMESTAMPTZ
);