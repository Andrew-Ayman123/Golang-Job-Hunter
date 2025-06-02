-- +goose Up
CREATE TABLE applicant_skills (
    applicant_id UUID NOT NULL REFERENCES applicants(user_id) ON DELETE CASCADE,
    skill_id INT NOT NULL REFERENCES skills(id) ON DELETE CASCADE,
    PRIMARY KEY (applicant_id, skill_id)
);

-- +goose Down
DROP TABLE IF EXISTS applicant_skills;
