-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    full_name TEXT NOT NULL,
    location TEXT,
    title TEXT,
    about_section TEXT,
    profile_picture TEXT,
    role TEXT NOT NULL CHECK (role IN ('applicant', 'recruiter','admin')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE TABLE applicants (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    resume TEXT
);

CREATE TABLE companies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE TABLE recruiters (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE SET NULL
);

CREATE TABLE admins (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    admin_level INT NOT NULL CHECK (admin_level >= 1)
);

-- +goose Down
DROP TABLE IF EXISTS admins;
DROP TABLE IF EXISTS recruiters;
DROP TABLE IF EXISTS companies;
DROP TABLE IF EXISTS applicants;
DROP TABLE IF EXISTS users;