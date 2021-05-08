drop table if exists go_test.anime;
create table anime (
    id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
    title VARCHAR(128) NOT NULL,
    content TEXT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO go_test.anime (title, content) VALUES ('天元突破グレンラガン', 'ガイナックス');
INSERT INTO go_test.anime (title) VALUES ('ガールズアンドパンツァー');
