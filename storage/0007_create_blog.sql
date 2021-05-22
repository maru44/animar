DROP TABLE IF EXISTS go_test.tbl_blog;
create table tbl_blog (
    id INT unsigned AUTO_INCREMENT NOT NULL PRIMARY KEY,
    content TEXT NULL,
    star TINYINT NULL,
    anime_ids INT unsigned NULL,
    user_id VARCHAR(128) NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4;
-- relation MM
DROP TABLE IF EXISTS go_test.relation_blog_animes;
CREATE TABLE relation_blog_animes (
    blog_id INT UNSIGNED NOT NULL,
    anime_id INT UNSIGNED NOT NULL,
    PRIMARY KEY (blog_id, anime_id)
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4;
ALTER TABLE relation_blog_animes
ADD CONSTRAINT fk_rel_blog_anime_id FOREIGN KEY (anime_id) REFERENCES anime (id) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE relation_blog_animes
ADD CONSTRAINT fk_rel_blog_blog_id FOREIGN KEY (blog_id) REFERENCES tbl_blog (id) ON DELETE CASCADE ON UPDATE CASCADE;