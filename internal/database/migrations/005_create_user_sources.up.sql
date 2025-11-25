CREATE TABLE user_sources (
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    source_id INTEGER REFERENCES sources(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, source_id)
);

CREATE INDEX idx_user_sources_user ON user_sources(user_id);

CREATE INDEX idx_user_sources_source ON user_sources(source_id);