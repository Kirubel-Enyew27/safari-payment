CREATE EXTENSION IF NOT EXISTS "uuid-ossp";  -- Enables uuid_generate_v4()

CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    checkout_request_id VARCHAR(100) UNIQUE NOT NULL,
    merchant_request_id VARCHAR(100) UNIQUE NOT NULL,
    phone_number VARCHAR(50) NOT NULL,
    amount NUMERIC(10, 2) NOT NULL,
    mpesa_receipt TEXT NOT NULL,
    transaction_date TIMESTAMP NOT NULL,
    result_code INT NOT NULL,
    result_desc TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
