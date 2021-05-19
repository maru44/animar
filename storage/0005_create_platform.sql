DROP TABLE IF EXISTS go_test.platform;
CREATE TABLE platform (
    id INT unsigned AUTO_INCREMENT NOT NULL PRIMARY KEY,
    eng_name VARCHAR(128) NOT NULL,
    plat_name VARCHAR(128) NULL,
    base_url VARCHAR(64) NULL,
    image VARCHAR(128) NULL,
    is_valid BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4;
-- add dummy
INSERT INTO go_test.platform (eng_name, plat_name)
VALUES ("netflix", "Netflix");
INSERT INTO go_test.platform (eng_name, plat_name)
VALUES ("primeVideo", "Prime video");
INSERT INTO go_test.platform (eng_name, plat_name)
VALUES ("dAnime", "dアニメ");