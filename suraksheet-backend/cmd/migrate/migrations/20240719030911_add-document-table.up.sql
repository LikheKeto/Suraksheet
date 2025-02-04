CREATE TABLE IF NOT EXISTS documents (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    referenceName VARCHAR(255) NOT NULL,
    bin INT NOT NULL,
    url TEXT NOT NULL,
    extract VARCHAR(1000) NOT NULL DEFAULT '',
    createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    language VARCHAR(3) NOT NULL DEFAULT 'eng',

    FOREIGN KEY (bin) REFERENCES bins(id) ON DELETE CASCADE,
    UNIQUE(referenceName, bin)
);