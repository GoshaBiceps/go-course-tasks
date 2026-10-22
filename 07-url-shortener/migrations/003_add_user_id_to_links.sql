ALTER TABLE links
ADD COLUMN user_id BIGINT;

ALTER TABLE links
ADD CONSTRAINT links_user_id_fkey
FOREIGN KEY (user_id)
REFERENCES users(id);