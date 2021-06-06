drop table if exists go_test.reviews;
create table reviews (
    id INT unsigned AUTO_INCREMENT NOT NULL PRIMARY KEY,
    content VARCHAR(160) NULL,
    rating TINYINT NULL,
    anime_id INT unsigned NOT NULL,
    user_id VARCHAR(128) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4;
-- foreign key
ALTER TABLE reviews
ADD CONSTRAINT fk_anime_reviews FOREIGN KEY (anime_id) REFERENCES animes (id) ON DELETE CASCADE ON UPDATE CASCADE;
-- index
ALTER TABLE reviews
ADD INDEX review_anime (anime_id);
ALTER TABLE reviews
ADD INDEX review_user (user_id);
-- add dummy
INSERT INTO go_test.reviews (content, rating, anime_id)
VALUES ('戦車\n国柄デフォルメ', 5, 2);