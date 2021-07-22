ALTER TABLE animes DROP COLUMN official_url;
ALTER TABLE animes DROP COLUMN twitter_url;
ALTER TABLE animes DROP COLUMN hash_tag;
ALTER TABLE animes
ADD COLUMN official_url VARCHAR(255) NULL
AFTER count_episodes;
ALTER TABLE animes
ADD COLUMN twitter_url VARCHAR(255) NULL
AFTER count_episodes;
ALTER TABLE animes
ADD COLUMN hash_tag VARCHAR(255) NULL
AFTER count_episodes;