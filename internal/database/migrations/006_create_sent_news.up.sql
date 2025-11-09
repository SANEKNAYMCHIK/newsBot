CREATE TABLE sent_news (
    news_id INTEGER REFERENCES news_items(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (news_id, user_id)
);