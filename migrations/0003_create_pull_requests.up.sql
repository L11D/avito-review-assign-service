CREATE TABLE pull_requests (
    id VARCHAR(50) PRIMARY KEY,
    name TEXT NOT NULL,
    status text NOT NULL DEFAULT 'OPEN',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    merged_at TIMESTAMP WITH TIME ZONE,
    author_id VARCHAR(50) NOT NULL REFERENCES users(id)
);