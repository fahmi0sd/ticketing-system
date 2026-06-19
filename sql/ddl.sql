-- Create table users 
CREATE TABLE IF NOT EXISTS users (
    id         SERIAL PRIMARY KEY,
    full_name  VARCHAR(255)        NOT NULL,
    email      VARCHAR(255) UNIQUE NOT NULL,
    password   VARCHAR(255)        NOT NULL,
    phone      VARCHAR(20),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create table routes
CREATE TABLE IF NOT EXISTS routes (
    id           SERIAL PRIMARY KEY,
    origin       VARCHAR(100)   NOT NULL,            
    destination  VARCHAR(100)   NOT NULL,            
    type         VARCHAR(20)    NOT NULL             
                 CHECK (type IN ('bus', 'train', 'flight')),
    operator     VARCHAR(100)   NOT NULL,            
    departure_at TIMESTAMP      NOT NULL,
    arrival_at   TIMESTAMP      NOT NULL,
    price        NUMERIC(12, 2) NOT NULL,
    quota        INT            NOT NULL DEFAULT 50,
    sold         INT            NOT NULL DEFAULT 0,  
    created_at   TIMESTAMP      DEFAULT NOW(),
    updated_at   TIMESTAMP      DEFAULT NOW()
);

-- Create table bookings
CREATE TABLE IF NOT EXISTS bookings (
    id          SERIAL PRIMARY KEY,
    user_id     INT            NOT NULL REFERENCES users(id),
    route_id    INT            NOT NULL REFERENCES routes(id),
    quantity    INT            NOT NULL DEFAULT 1,
    total_price NUMERIC(12, 2) NOT NULL,
    status      VARCHAR(20)    NOT NULL DEFAULT 'pending'
                CHECK (status IN ('pending', 'paid', 'cancelled', 'expired')),
    payment_url TEXT,                                
    external_id VARCHAR(100),                       
    expired_at  TIMESTAMP,                           
    paid_at     TIMESTAMP,
    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP DEFAULT NOW()
);

-- Insert seed data route
INSERT INTO routes (origin, destination, type, operator, departure_at, arrival_at, price, quota)
VALUES
    ('Surabaya', 'Makassar', 'flight', 'Lion Air', NOW() + INTERVAL '2 days 7 hours', NOW() + INTERVAL '2 days 8 hours 30 minutes', 950000, 150),
    ('Surabaya', 'Balikpapan', 'flight', 'Citilink', NOW() + INTERVAL '3 days 10 hours', NOW() + INTERVAL '3 days 11 hours 15 minutes', 1100000, 120),
    ('Surabaya', 'Jakarta', 'flight', 'Batik Air', NOW() + INTERVAL '1 day 6 hours', NOW() + INTERVAL '1 day 7 hours 30 minutes', 1200000, 140),
    ('Malang', 'Jakarta', 'flight', 'Garuda Indonesia', NOW() + INTERVAL '4 days 8 hours', NOW() + INTERVAL '4 days 9 hours 30 minutes', 1600000, 96),
    ('Surabaya', 'Malang', 'train', 'KAI (Lokal Penataran)', NOW() + INTERVAL '1 day 2 hours', NOW() + INTERVAL '1 day 4 hours', 15000, 80),
    ('Surabaya', 'Yogyakarta', 'train', 'KAI (Sancaka)', NOW() + INTERVAL '2 days 3 hours', NOW() + INTERVAL '2 days 7 hours', 240000, 64),
    ('Banyuwangi', 'Surabaya', 'train', 'KAI (Probowangi)', NOW() + INTERVAL '1 day 5 hours', NOW() + INTERVAL '1 day 11 hours', 56000, 100),
    ('Malang', 'Jakarta', 'train', 'KAI (Gajayana)', NOW() + INTERVAL '3 days 16 hours', NOW() + INTERVAL '4 days 4 hours', 650000, 50),
    ('Madiun', 'Surabaya', 'train', 'KAI (Jayakarta)', NOW() + INTERVAL '2 days 13 hours', NOW() + INTERVAL '2 days 15 hours 30 minutes', 180000, 70),
    ('Surabaya', 'Malang', 'bus', 'Restu Panda', NOW() + INTERVAL '1 day 1 hour', NOW() + INTERVAL '1 day 2 hours 30 minutes', 40000, 45),
    ('Surabaya', 'Ponorogo', 'bus', 'Harapan Jaya', NOW() + INTERVAL '1 day 8 hours', NOW() + INTERVAL '1 day 12 hours', 70000, 40),
    ('Jember', 'Surabaya', 'bus', 'Akas Asri', NOW() + INTERVAL '2 days 4 hours', NOW() + INTERVAL '2 days 8 hours 30 minutes', 85000, 45),
    ('Surabaya', 'Jakarta', 'bus', 'Rosalia Indah', NOW() + INTERVAL '3 days 19 hours', NOW() + INTERVAL '4 days 5 hours', 375000, 32),
    ('Surabaya', 'Makassar', 'flight', 'Citilink', NOW() + INTERVAL '2 days 7 hours', NOW() + INTERVAL '2 days 8 hours 30 minutes', 1100000, 120),
    ('Surabaya', 'Balikpapan', 'flight', 'Super Air Jet', NOW() + INTERVAL '3 days 10 hours', NOW() + INTERVAL '3 days 11 hours 15 minutes', 1050000, 180),
    ('Surabaya', 'Jakarta', 'flight', 'Garuda Indonesia', NOW() + INTERVAL '1 day 6 hours', NOW() + INTERVAL '1 day 7 hours 30 minutes', 1650000, 96),
    ('Malang', 'Jakarta', 'flight', 'Batik Air', NOW() + INTERVAL '4 days 8 hours', NOW() + INTERVAL '4 days 9 hours 30 minutes', 1250000, 140),
    ('Surabaya', 'Malang', 'train', 'KAI (Arjuno Ekspres)', NOW() + INTERVAL '1 day 2 hours', NOW() + INTERVAL '1 day 4 hours', 75000, 50),
    ('Surabaya', 'Yogyakarta', 'train', 'KAI (Logawa)', NOW() + INTERVAL '2 days 3 hours', NOW() + INTERVAL '2 days 7 hours 30 minutes', 190000, 80),
    ('Banyuwangi', 'Surabaya', 'train', 'KAI (Sri Tanjung)', NOW() + INTERVAL '1 day 5 hours', NOW() + INTERVAL '1 day 11 hours 15 minutes', 94000, 90),
    ('Malang', 'Jakarta', 'train', 'KAI (Brawijaya)', NOW() + INTERVAL '3 days 16 hours', NOW() + INTERVAL '4 days 4 hours 10 minutes', 580000, 60),
    ('Surabaya', 'Malang', 'bus', 'PO Tentrem', NOW() + INTERVAL '1 day 1 hour', NOW() + INTERVAL '1 day 2 hours 30 minutes', 40000, 45),
    ('Surabaya', 'Ponorogo', 'bus', 'PO Jaya Utama', NOW() + INTERVAL '1 day 8 hours', NOW() + INTERVAL '1 day 12 hours', 65000, 50),
    ('Jember', 'Surabaya', 'bus', 'PO Ladju', NOW() + INTERVAL '2 days 4 hours', NOW() + INTERVAL '2 days 8 hours 30 minutes', 85000, 45),
    ('Surabaya', 'Jakarta', 'bus', 'PO Sinar Jaya', NOW() + INTERVAL '3 days 19 hours', NOW() + INTERVAL '4 days 5 hours', 350000, 40);

