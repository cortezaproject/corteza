CREATE TABLE IF NOT EXISTS sys_reminder (
    `id`           BIGINT(20)   UNSIGNED NOT NULL,
    `resource`     VARCHAR(128)          NOT NULL                           COMMENT 'Resource, that this reminder is bound to',
    `payload`      JSON                  NOT NULL                           COMMENT 'Payload for this reminder',
    `snooze_count` INT                   NOT NULL DEFAULT 0                 COMMENT 'Number of times this reminder was snoozed',

    `assigned_to`  BIGINT(20)   UNSIGNED NOT NULL DEFAULT 0                 COMMENT 'Assignee for this reminder',
    `assigned_by`  BIGINT(20)   UNSIGNED NOT NULL DEFAULT 0                 COMMENT 'User that assigned this reminder',
    `assigned_at`  DATETIME              NOT NULL                           COMMENT 'When the reminder was assigned',

    `dismissed_by` BIGINT(20)   UNSIGNED NOT NULL DEFAULT 0                 COMMENT 'User that dismissed this reminder',
    `dismissed_at` DATETIME                  NULL DEFAULT NULL              COMMENT 'Time the reminder was dismissed',

    `remind_at`    DATETIME                  NULL DEFAULT NULL              COMMENT 'Time the user should be reminded',

    `created_by`   BIGINT(20)  UNSIGNED NOT NULL DEFAULT 0,
    `created_at`   DATETIME             NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_by`   BIGINT(20)  UNSIGNED NOT NULL DEFAULT 0,
    `updated_at`   DATETIME                 NULL DEFAULT NULL,
    `deleted_by`   BIGINT(20)  UNSIGNED NOT NULL DEFAULT 0,
    `deleted_at`   DATETIME                 NULL DEFAULT NULL,

    PRIMARY KEY (`id`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8;
