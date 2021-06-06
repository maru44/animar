drop table if exists go_test.reviews;
drop table if exists go_test.watch_states;
drop table if exists go_test.animes;
drop table if exists go_test.relation_anime_season;
create table animes (
    id INT unsigned AUTO_INCREMENT NOT NULL PRIMARY KEY,
    slug CHAR(12) NOT NULL UNIQUE,
    title VARCHAR(128) NOT NULL,
    abbreviation VARCHAR(64) NULL COMMENT '略称',
    kana VARCHAR(64) NULL COMMENT 'カナ',
    eng_name VARCHAR(64) COMMENT 'roman',
    thumb_url VARCHAR(256) NULL,
    copyright VARCHAR(256) NULL,
    description TEXT NULL,
    -- on_air_state TINYINT NULL COMMENT '0: 打ち切り, 1: 放送終了, 2: 放送中, 3: 放送前',
    state VARCHAR(16) NULL COMMENT 'cut: 打ち切り, fin: 放送終了, now: 放送中, pre: 放送前',
    series_id INT unsigned NULL,
    count_episodes INT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP -- INDEX anime_slug (slug)
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4;
-- index
ALTER TABLE animes
ADD UNIQUE INDEX anime_slug (slug);
ALTER TABLE animes
ADD UNIQUE INDEX anime_series (series_id);
-- foreign key
ALTER TABLE animes
ADD CONSTRAINT fk_anime_series FOREIGN KEY (series_id) REFERENCES series (id) ON DELETE
SET NULL ON UPDATE CASCADE;
-- MM relation (season)
-- relation
CREATE TABLE relation_anime_season (
    season_id INT UNSIGNED NOT NULL,
    anime_id INT UNSIGNED NOT NULL,
    PRIMARY KEY (season_id, anime_id)
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4;
ALTER TABLE relation_season_anime
ADD CONSTRAINT fk_relation_season_anime_id FOREIGN KEY (anime_id) REFERENCES animes (id) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE relation_season_season
ADD CONSTRAINT fk_relation_season_season_id FOREIGN KEY (season_id) REFERENCES season (id) ON DELETE CASCADE ON UPDATE CASCADE;
-- add dummy
INSERT INTO go_test.animes (title, slug, description)
VALUES ('天元突破グレンラガン', 'aaaaaaaaaaaa', 'ガイナックス');
INSERT INTO go_test.animes (title, slug, series_id)
VALUES ('ガールズアンドパンツァー', 'aaaaaaaaaaab', 1);
INSERT INTO go_test.animes (title, slug)
VALUES ('宇宙よりも遠い場所', 'aaaaaaaaaaac');
INSERT INTO go_test.animes (title, slug)
VALUES ('ゾンビランドサガ', 'aaaaaaaaaaad');
INSERT INTO go_test.animes (title, slug, description)
VALUES ('あしたのジョー', 'aaaaaaaaaaae', '明日はどっちだ?');