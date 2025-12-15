CREATE TABLE sources (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    url VARCHAR(500) UNIQUE NOT NULL,
    category_id INTEGER REFERENCES categories(id),
    is_active BOOLEAN DEFAULT FALSE
);

CREATE INDEX idx_sources_category ON sources(category_id);
CREATE INDEX idx_sources_checked ON sources(is_active) WHERE is_active = true;
