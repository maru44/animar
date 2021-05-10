drop table if exists go_test.tbl_reviews;
create table tbl_reviews (
    id INT unsigned AUTO_INCREMENT NOT NULL PRIMARY KEY,
    content TEXT NULL,
    star INT NULL,
    anime_id INT unsigned NOT NULL,
    user_id VARCHAR(128) NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=INNODB DEFAULT CHARSET=utf8mb4;

ALTER TABLE tbl_reviews
    ADD CONSTRAINT fk_anime_reviews
    FOREIGN KEY (anime_id)
    REFERENCES anime (id)
    ON DELETE CASCADE ON UPDATE CASCADE;

INSERT INTO go_test.tbl_reviews (content, star, anime_id) VALUES ('戦車\n国柄デフォルメ', 5, 2);
INSERT INTO go_test.tbl_reviews (content, user_id, anime_id) VALUES ('iine', 'NiimLADqi5RJaRqCdatNFgKLKnK2', 2);
