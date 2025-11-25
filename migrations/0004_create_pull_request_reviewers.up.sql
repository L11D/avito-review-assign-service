CREATE TABLE pull_request_reviewers (
    user_id VARCHAR(50) NOT NULL REFERENCES users(id),
    pull_request_id VARCHAR(50) NOT NULL REFERENCES pull_requests(id),
    PRIMARY KEY (user_id, pull_request_id)
);