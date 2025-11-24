CREATE TABLE teams (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name varchar(50) NOT NULL UNIQUE
);

INSERT INTO teams (name) VALUES ('Engineering');
INSERT INTO teams (name) VALUES ('Marketing');