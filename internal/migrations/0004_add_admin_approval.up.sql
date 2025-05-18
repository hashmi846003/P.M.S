-- Add admin approval columns (if not included in initial schema)
ALTER TABLE users
ADD COLUMN status VARCHAR(20) NOT NULL DEFAULT 'pending',
ADD COLUMN role VARCHAR(20) NOT NULL DEFAULT 'user';