-- Create page versions table
CREATE TABLE page_versions (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    page_id INT NOT NULL REFERENCES pages(id),
    content TEXT NOT NULL
);

-- Create indexes
CREATE INDEX idx_page_versions_page_id ON page_versions(page_id);