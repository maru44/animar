drop table if exists go_test.anime_platform;
create table anime_platform (
    platform_id INT unsigned NOT NULL,
    anime_id INT unsigned NOT NULL,
    user_id VARCHAR(128) NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(platform_id, anime_id)
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4;
ALTER TABLE anime_platform
ADD CONSTRAINT fk_anime_platform_anime_id FOREIGN KEY (anime_id) REFERENCES anime (id) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE anime_platform
ADD CONSTRAINT fk_anime_platform_platform_id FOREIGN KEY (platform_id) REFERENCES platform (id) ON DELETE CASCADE ON UPDATE CASCADE;
-- add dummy
INSERT INTO go_test.anime_platform (platform_id, anime_id)
VALUES (3, 5);