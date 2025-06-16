-- +goose Up
CREATE TABLE skills (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    is_default BOOLEAN DEFAULT false
);

INSERT INTO skills (name, is_default) VALUES
    ('JavaScript', true),
    ('Python', true),
    ('Java', true),
    ('C#', true),
    ('Go', true),
    ('Ruby', true),
    ('SQL', true),
    ('HTML', true),
    ('CSS', true),
    ('React', true),
    ('Vue.js', true),
    ('Node.js', true),
    ('Django', true),
    ('Flask', true),
    ('Spring Boot', true),
    ('Project Management', true),
    ('Agile Methodology', true),
    ('UI/UX Design', true),
    ('Marketing', true),
    ('Sales', true),
    ('DevOps', true),
    ('Cloud Computing', true),
    ('AWS', true),
    ('Azure', true),
    ('Machine Learning', true),
    ('Data Analysis', true),
    ('Cybersecurity', true);

-- +goose Down
DROP TABLE IF EXISTS skills;
