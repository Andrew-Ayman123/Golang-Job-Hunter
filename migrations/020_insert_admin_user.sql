-- +goose Up
INSERT INTO users (email, password_hash,full_name, role)
VALUES (
    'andrewayman9@gmail.com',
    '$2a$10$S9sO2OMJI0Q.uGrsxx/wRu4q7QEcNP8XcfEBvFYGkrm1fwlykDqfC', -- "12345678"
    'Andrew Ayman',
    'admin'
);

-- Insert into admins using the ID of the new user
INSERT INTO admins (user_id, admin_level)
SELECT id, 1 FROM users WHERE email = 'andrewayman9@gmail.com';

-- +goose Down
DELETE FROM admins WHERE user_id = (SELECT id FROM users WHERE email = 'andrewayman9@gmail.com');
DELETE FROM users WHERE email = 'andrewayman9@gmail.com';
