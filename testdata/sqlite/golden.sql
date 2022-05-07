PRAGMA foreign_keys = false;

DROP TABLE IF EXISTS `player`;

CREATE TABLE `player` (
    `id` INTEGER NOT NULL,
    `name` TEXT NOT NULL,
    `created_at` INTEGER NOT NULL,
    `updated_at` INTEGER NOT NULL,
    `daily_notification_at` INTEGER NOT NULL,
    PRIMARY KEY (`id`)
);


DROP TABLE IF EXISTS `entry`;

CREATE TABLE `entry` (
    `id` INTEGER NOT NULL,
    `player_id` INTEGER NOT NULL,
    `title` TEXT NOT NULL,
    `public` INTEGER NOT NULL DEFAULT 0,
    `content` TEXT NOT NULL,
    `created_at` INTEGER NOT NULL,
    `updated_at` INTEGER NOT NULL,
    FOREIGN KEY (`player_id`) REFERENCES `player` (`id`) ON DELETE CASCADE,
    PRIMARY KEY (`id`)
);

CREATE INDEX `created_at_idx` ON `entry` (`created_at`);
CREATE INDEX `title_idx` ON `entry` (`title`);
CREATE UNIQUE INDEX `created_at_uniq_idx` ON `entry` (`created_at`);
PRAGMA foreign_keys = true;
