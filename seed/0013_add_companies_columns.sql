ALTER TABLE companies
ADD COLUMN explanation TEXT NULL
AFTER official_url;
ALTER TABLE companies
ADD COLUMN twitter_account VARCHAR(128) NULL
AFTER official_url;