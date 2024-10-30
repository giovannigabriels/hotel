CREATE TABLE IF NOT EXISTS refunds (
    id SERIAL PRIMARY KEY,
    payment_id INT NOT NULL REFERENCES payments(id) ON DELETE CASCADE,
    refund_amount DECIMAL(10, 2) NOT NULL,
    refund_status VARCHAR(20) DEFAULT 'requested', -- Options: 'requested', 'completed', 'denied'
    refund_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);