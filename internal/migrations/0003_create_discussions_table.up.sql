-- Create discussions table
CREATE TABLE discussions (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    page_id INT NOT NULL REFERENCES pages(id),
    user_id INT NOT NULL REFERENCES users(id),
    content TEXT NOT NULL
);

-- Create indexes
CREATE INDEX idx_discussions_page_id ON discussions(page_id);
CREATE INDEX idx_discussions_user_id ON discussions(user_id);