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
    duration INT NOT NULL,   -- тривалість зарядки у хвилинах
    energy INT NOT NULL,     -- спожита енергія у кВт·год
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
