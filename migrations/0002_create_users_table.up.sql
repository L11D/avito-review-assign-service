CREATE TABLE users (
    id VARCHAR(50) PRIMARY KEY,
    username TEXT NOT NULL,
    is_active BOOLEAN NOT NULL,
    team_id UUID NOT NULL REFERENCES teams(id)
);