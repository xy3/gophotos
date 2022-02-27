PRAGMA journal_mode = MEMORY;
PRAGMA synchronous = OFF;
PRAGMA foreign_keys = OFF;
PRAGMA ignore_check_constraints = OFF;
PRAGMA auto_vacuum = NONE;
PRAGMA secure_delete = OFF;
BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS `photos`
(
    `id`               INTEGER  PRIMARY KEY AUTOINCREMENT,
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL     DEFAULT NULL,
    `size`       INTEGER   NOT NULL,
    `file_name`  TEXT      NOT NULL,
    `file_hash`  TEXT      NOT NULL,
    `extension`  TEXT      NOT NULL,
    `user_id`    INTEGER   NOT NULL,
    UNIQUE(`user_id`, `file_hash`)
);

CREATE TABLE IF NOT EXISTS `users`
(
    `id`               INTEGER  PRIMARY KEY AUTOINCREMENT,
    `created_at`       timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`       timestamp NULL     DEFAULT NULL,
    `email`            TEXT      NOT NULL,
    `password`         TEXT      NOT NULL,
    `storage_path`     TEXT               DEFAULT NULL,
    UNIQUE (`email`)
);

-- INSERT into `photos`
-- Values (1, 0, null, 10, 'fileNameTest', 'fileHashTest', 'jpg', 123);
COMMIT;
PRAGMA ignore_check_constraints = ON;
PRAGMA foreign_keys = ON;
PRAGMA journal_mode = WAL;
PRAGMA synchronous = NORMAL;
