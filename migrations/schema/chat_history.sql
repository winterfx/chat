CREATE TABLE chat_history (
    id INT AUTO_INCREMENT PRIMARY KEY,
    message_id CHAR(36) NOT NULL UNIQUE,  -- UUID
    message_timestamp TIMESTAMP NOT NULL,
    sender_id VARCHAR(255) NOT NULL,
    receiver_id VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    retry_count INT DEFAULT 0,
    status VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);