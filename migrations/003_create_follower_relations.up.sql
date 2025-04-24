CREATE TABLE follower_relations (
                                    id SERIAL PRIMARY KEY,
                                    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                    follower_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                    status VARCHAR(20) DEFAULT 'pending',
                                    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Уникальная пара подписки
CREATE UNIQUE INDEX uniq_follower_pair ON follower_relations(user_id, follower_id);

-- Индексы для быстрого поиска
CREATE INDEX idx_follower_user_id ON follower_relations(user_id);
CREATE INDEX idx_follower_follower_id ON follower_relations(follower_id);