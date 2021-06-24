DROP TABLE IF EXISTS go_test.relation_blog_animes;
DROP TABLE IF EXISTS go_test.blogs;
create table blogs (
    id INT unsigned AUTO_INCREMENT NOT NULL PRIMARY KEY,
    slug VARCHAR(16) NOT NULL,
    title VARCHAR(64) NOT NULL,
    abstract VARCHAR(160) NULL,
    content TEXT NOT NULL,
    user_id VARCHAR(128) NOT NULL,
    is_public BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4;
-- index
ALTER TABLE blogs
ADD INDEX blog_user (user_id);
ALTER TABLE blogs
ADD INDEX blog_user_public (user_id, is_public);
ALTER TABLE blogs
ADD UNIQUE INDEX blog_slug (slug);
-- relation MM
CREATE TABLE relation_blog_animes (
    blog_id INT UNSIGNED NOT NULL,
    anime_id INT UNSIGNED NOT NULL,
    PRIMARY KEY (blog_id, anime_id)
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4;
ALTER TABLE relation_blog_animes
ADD CONSTRAINT fk_rel_blog_anime_id FOREIGN KEY (anime_id) REFERENCES animes (id) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE relation_blog_animes
ADD CONSTRAINT fk_rel_blog_blog_id FOREIGN KEY (blog_id) REFERENCES blogs (id) ON DELETE CASCADE ON UPDATE CASCADE;