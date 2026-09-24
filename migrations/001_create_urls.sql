CREATE TABLE urls (
    id SERIAL PRIMARY KEY,
    short_code VARCHAR(10) UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    access_count INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);