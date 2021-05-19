drop table if exists go_test.tbl_reviews;
drop table if exists go_test.watch_states;
drop table if exists go_test.anime;
create table anime (
    id INT unsigned AUTO_INCREMENT NOT NULL PRIMARY KEY,
    slug CHAR(12) NOT NULL UNIQUE,
    title VARCHAR(128) NOT NULL,
    -- abbreviation VARCHAR(64) NULL,
    content TEXT NULL,
    on_air_state TINYINT NULL COMMENT '0: 打ち切り, 1: 放送終了, 2: 放送中, 3: 放送前',
    series_id INT unsigned NULL,
    season VARCHAR(256) NULL,
    stories INT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP -- INDEX anime_slug (slug)
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4;
ALTER TABLE anime
ADD CONSTRAINT fk_anime_series FOREIGN KEY (series_id) REFERENCES series (id) ON DELETE
SET NULL ON UPDATE CASCADE;
-- add dummy
INSERT INTO go_test.anime (title, slug, content)
VALUES ('天元突破グレンラガン', 'aaaaaaaaaaaa', 'ガイナックス');
INSERT INTO go_test.anime (title, slug, series_id)
VALUES ('ガールズアンドパンツァー', 'aaaaaaaaaaab', 1);
INSERT INTO go_test.anime (title, slug)
VALUES ('宇宙よりも遠い場所', 'aaaaaaaaaaac');
INSERT INTO go_test.anime (title, slug)
VALUES ('ゾンビランドサガ', 'aaaaaaaaaaad');
INSERT INTO go_test.anime (title, slug, content)
VALUES ('あしたのジョー', 'aaaaaaaaaaae', '明日はどっちだ?');