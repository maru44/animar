DROP TABLE IF EXISTS go_test.series;
CREATE TABLE series (
    id INT unsigned AUTO_INCREMENT NOT NULL PRIMARY KEY,
    eng_name VARCHAR(128) NOT NULL,
    series_name VARCHAR(128) NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=INNODB DEFAULT CHARSET=utf8mb4;

INSERT INTO go_test.series (eng_name, series_name) VALUES ("girls_und_panzer", "ガールズアンドパンツァー");
