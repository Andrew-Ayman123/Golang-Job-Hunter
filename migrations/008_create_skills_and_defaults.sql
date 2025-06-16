-- +goose Up
CREATE TABLE skills (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

INSERT INTO skills (name) VALUES
    ('JavaScript'),
    ('Python'),
    ('Java'),
    ('C#'),
    ('Go'),
    ('Ruby'),
    ('SQL'),
    ('HTML'),
    ('CSS'),
    ('React'),
    ('Vue.js'),
    ('Node.js'),
    ('Django'),
    ('Flask'),
    ('Spring Boot'),
    ('Project Management'),
    ('Agile Methodology'),
    ('UI/UX Design'),
    ('Marketing'),
    ('Sales'),
    ('DevOps'),
    ('Cloud Computing'),
    ('AWS'),
    ('Azure'),
    ('Machine Learning'),
    ('Data Analysis'),
    ('Cybersecurity');

-- +goose Down
DROP TABLE IF EXISTS skills;
