ALTER TABLE users
ADD COLUMN followers_count INT DEFAULT 0,
ADD COLUMN following_count INT DEFAULT 0;