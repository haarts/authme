CREATE TABLE `users` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `username` VARCHAR NOT NULL UNIQUE,
    `encrypted_password` VARCHAR NOT NULL,
    `salt` VARCHAR NOT NULL,
    `email` VARCHAR NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE `sessions` (
    `session` VARCHAR NOT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP
);
