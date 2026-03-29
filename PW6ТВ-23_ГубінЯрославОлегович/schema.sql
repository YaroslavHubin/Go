CREATE TABLE stations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    power INT NOT NULL,
    slots INT NOT NULL,
    location VARCHAR(150)
);

CREATE TABLE sessions (
    id SERIAL PRIMARY KEY,
    station_id INT REFERENCES stations(id),
    duration INT NOT NULL,
    energy INT NOT NULL,
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(150) UNIQUE NOT NULL,
    password VARCHAR(200) NOT NULL,
    role VARCHAR(50) DEFAULT 'user'
);

-- Адміністраторський акаунт
INSERT INTO users (name, email, password, role)
VALUES ('admin', 'admin@local', 'admin', 'admin');
