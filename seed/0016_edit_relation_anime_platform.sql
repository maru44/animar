-- add interval relation_anime_platform
ALTER TABLE relation_anime_platform
ADD delivery_interval VARCHAR(32) NULL
AFTER link_url;
-- add first 
ALTER TABLE relation_anime_platform
ADD first_broadcast TIMESTAMP NULL DEFAULT NULL
AFTER delivery_interval;