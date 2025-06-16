-- +goose Up

-- Create phone numbers table
CREATE TABLE recruiters (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE SET NULL
);

CREATE TABLE admins (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    admin_level INT NOT NULL CHECK (admin_level >= 1)
);

-- +goose Down

DROP TABLE IF EXISTS recruiters;
DROP TABLE IF EXISTS admins;