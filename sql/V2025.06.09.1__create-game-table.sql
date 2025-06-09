CREATE TABLE games (
    id INT PRIMARY KEY NOT NULL,
    shortlink VARCHAR(255) NOT NULL UNIQUE,
    settings JSON
);