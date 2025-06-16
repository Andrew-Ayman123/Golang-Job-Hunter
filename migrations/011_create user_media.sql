-- +goose Up

-- Create media table for education, experience, certifications, and projects
CREATE TABLE user_media (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    media_type TEXT NOT NULL CHECK (media_type IN ('image', 'video', 'document')),
    file_name TEXT NOT NULL,
    file_path TEXT NOT NULL,
    file_size BIGINT,
    mime_type TEXT,
    alt_text TEXT,
    description TEXT,
    -- Reference fields to link media to specific entities
    education_id UUID REFERENCES user_education(id) ON DELETE CASCADE,
    experience_id UUID REFERENCES user_experience(id) ON DELETE CASCADE,
    certification_id UUID REFERENCES user_certifications(id) ON DELETE CASCADE,
    project_id UUID REFERENCES user_projects(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    -- Ensure media is linked to exactly one entity
    CONSTRAINT media_single_reference CHECK (
        (education_id IS NOT NULL)::int + 
        (experience_id IS NOT NULL)::int + 
        (certification_id IS NOT NULL)::int + 
        (project_id IS NOT NULL)::int = 1
    )
);

-- +goose Down

DROP TABLE IF EXISTS user_media;
