drop table if exists go_test.watch_states;
create table watch_states (
    id INT unsigned AUTO_INCREMENT NOT NULL PRIMARY KEY,
    watch TINYINT NULL COMMENT '0: 脱落, 1: 視聴中, 2: 完了, 3: 周回済み, 4: 興味',
    anime_id INT unsigned NOT NULL,
    user_id VARCHAR(128) NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=INNODB DEFAULT CHARSET=utf8mb4;

ALTER TABLE watch_states
    ADD CONSTRAINT fk_anime_watche_states
    FOREIGN KEY (anime_id)
    REFERENCES anime (id)
    ON DELETE CASCADE ON UPDATE CASCADE;

INSERT INTO go_test.watch_states (watch, anime_id, user_id) VALUES (3, 1, 'NiimLADqi5RJaRqCdatNFgKLKnK2');
INSERT INTO go_test.watch_states (watch, anime_id, user_id) VALUES (3, 2, 'sample1');
INSERT INTO go_test.watch_states (watch, anime_id, user_id) VALUES (2, 1, 'sample1');
INSERT INTO go_test.watch_states (watch, anime_id, user_id) VALUES (4, 1, 'sample2');
INSERT INTO go_test.watch_states (watch, anime_id, user_id) VALUES (1, 1, 'sample3');
