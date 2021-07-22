DROP TABLE IF EXISTS go_test.companies;
CREATE TABLE companies (
    id INT unsigned AUTO_INCREMENT NOT NULL PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    eng_name VARCHAR(128) NOT NULL,
    official_url VARCHAR(256) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX staff_eng_name (eng_name)
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4;
-- add animes column (company_id)
ALTER TABLE animes DROP COLUMN company_id;
ALTER TABLE animes
ADD COLUMN company_id INT UNSIGNED NULL
AFTER series_id;
ALTER TABLE animes
ADD CONSTRAINT fk_anime_company FOREIGN KEY (company_id) REFERENCES companies (id) ON DELETE
SET NULL ON UPDATE CASCADE;