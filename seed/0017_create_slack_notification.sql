DROP TABLE IF EXISTS go_test.slacks;
CREATE TABLE slacks (
    id INT unsigned AUTO_INCREMENT NOT NULL PRIMARY KEY,
    user_id VARCHAR(128),
    slack_channel VARCHAR(128),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX user_index(user_id),
    INDEX slack_index(slack_channel)
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4;