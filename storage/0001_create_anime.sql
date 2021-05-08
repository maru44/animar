drop table if exists go_test.anime;
create table anime (
    id INT unsigned AUTO_INCREMENT NOT NULL PRIMARY KEY,
    title VARCHAR(128) NOT NULL,
    content TEXT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=INNODB DEFAULT CHARSET=utf8mb4;

INSERT INTO go_test.anime (title, content) VALUES ('天元突破グレンラガン', 'ガイナックス');
INSERT INTO go_test.anime (title) VALUES ('ガールズアンドパンツァー');
