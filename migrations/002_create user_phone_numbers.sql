-- +goose Up

-- Create phone numbers table
CREATE TABLE user_phone_numbers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    phone_number TEXT NOT NULL,
    phone_type TEXT CHECK (phone_type IN ('mobile', 'home', 'work', 'other')) DEFAULT 'mobile',
    is_primary BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- +goose Down

DROP TABLE IF EXISTS user_phone_numbers;