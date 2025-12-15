CREATE TABLE news_items (
    id SERIAL PRIMARY KEY,
    title VARCHAR(500) NOT NULL,
    content TEXT,
    url VARCHAR(500) UNIQUE NOT NULL,
    published_at TIMESTAMPTZ NOT NULL,
    source_id INTEGER REFERENCES sources(id) ON DELETE CASCADE NOT NULL,
    guid VARCHAR(500) NOT NULL,
    UNIQUE(source_id, guid)
);

CREATE INDEX idx_news_source ON news_items(source_id);
CREATE INDEX idx_news_published ON news_items(published_at DESC);