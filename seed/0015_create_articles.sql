DROP TABLE IF EXISTS go_test.interview_quote;
DROP TABLE IF EXISTS go_test.article_chara;
DROP TABLE IF EXISTS go_test.articles;
DROP TABLE IF EXISTS go_test.relation_article_anime;
CREATE TABLE articles (
    id INT UNSIGNED AUTO_INCREMENT NOT NULL PRIMARY KEY,
    slug VARCHAR(16) NOT NULL,
    article_type VARCHAR(16) DEFAULT 'NEWS' COMMENT 'NEWS: ニュース, INTV: インタビュー',
    abstract TEXT NULL,
    content TEXT NULL,
    image VARCHAR(128) NULL,
    author VARCHAR(64) NULL,
    is_public BOOLEAN DEFAULT FALSE,
    user_id VARCHAR(128) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX slug_index(slug)
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4;
-- character of article
CREATE TABLE IF NOT EXISTS go_test.article_chara(
    id INT UNSIGNED AUTO_INCREMENT NOT NULL PRIMARY KEY,
    chara_name VARCHAR(48),
    image VARCHAR(128) NULL,
    user_id VARCHAR(128) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_article_chara FOREIGN KEY (article_id) REFERENCES articles (id) ON DELETE
    SET NULL ON UPDATE CASCADE
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4;
-- interview section
CREATE TABLE IF NOT EXISTS go_test.interview_quote(
    id INT UNSIGNED AUTO_INCREMENT NOT NULL PRIMARY KEY,
    article_id INT UNSIGNED NULL,
    chara_id INT UNSIGNED NULL,
    sequence INT DEFAULT 0,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_article_interview FOREIGN KEY (article_id) REFERENCES articles (id) ON DELETE
    SET NULL ON UPDATE CASCADE,
        CONSTRAINT fk_article_interviewer FOREIGN KEY (chara_id) REFERENCES article_chara (id) ON DELETE
    SET NULL ON UPDATE CASCADE
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4;
-- Relation article + character
CREATE TABLE IF NOT EXISTS go_test.relation_article_chara (
    article_id INT UNSIGNED NOT NULL,
    chara_id INT UNSIGNED NOT NULL,
    PRIMARY KEY(article_id, chara_id)
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4;
-- Relation anime + article(M/M)
CREATE TABLE IF NOT EXISTS go_test.relation_article_anime (
    anime_id INT UNSIGNED NOT NULL,
    article_id INT UNSIGNED NOT NULL,
    PRIMARY KEY(anime_id, article_id)
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4;