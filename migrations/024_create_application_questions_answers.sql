-- +goose Up
CREATE TABLE application_questions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id UUID NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    question TEXT NOT NULL,
    is_required BOOLEAN DEFAULT true,
    question_order INT DEFAULT 0
);

CREATE TABLE application_answers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id UUID NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    question_id UUID NOT NULL REFERENCES application_questions(id) ON DELETE CASCADE,
    answer TEXT,
    UNIQUE(application_id, question_id)
);

-- +goose Down
DROP TABLE IF EXISTS application_answers;
DROP TABLE IF EXISTS application_questions;
