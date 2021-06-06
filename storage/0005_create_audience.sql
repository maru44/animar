drop table if exists go_test.audiences;
create table audiences (
    id INT unsigned AUTO_INCREMENT NOT NULL PRIMARY KEY,
    state TINYINT NULL COMMENT '0: 脱落, 1: 興味, 2: 視聴中, 3: 完了, 4: 周回済み',
    anime_id INT unsigned NOT NULL,
    user_id VARCHAR(128) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4;
-- foreign key
ALTER TABLE audiences
ADD CONSTRAINT fk_anime_audience FOREIGN KEY (anime_id) REFERENCES anime (id) ON DELETE CASCADE ON UPDATE CASCADE;
-- index
ALTER TABLE audiences
ADD INDEX audience_anime (anime_id);
ALTER TABLE audiences
ADD INDEX audience_user (user_id);
-- add dummy
INSERT INTO go_test.audiences (state, anime_id, user_id)
VALUES (4, 1, 'sample4');
INSERT INTO go_test.audiences (state, anime_id, user_id)
VALUES (4, 2, 'sample1');
INSERT INTO go_test.audiences (state, anime_id, user_id)
VALUES (3, 1, 'sample1');
INSERT INTO go_test.audiences (state, anime_id, user_id)
VALUES (2, 1, 'sample2');
INSERT INTO go_test.audiences (state, anime_id, user_id)
VALUES (1, 1, 'sample3');