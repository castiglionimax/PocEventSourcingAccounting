CREATE TABLE IF NOT EXISTS accounts (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    amount DECIMAL(50, 6) NOT NULL,
    number VARCHAR(50) NOT NULL,
    last_updated DATETIME NOT NULL
    );

