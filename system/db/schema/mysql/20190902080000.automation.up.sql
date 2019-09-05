CREATE TABLE IF NOT EXISTS sys_automation_script (
    `id`            BIGINT(20)  UNSIGNED NOT NULL,
    `rel_namespace` BIGINT(20)  UNSIGNED NOT NULL DEFAULT 0         COMMENT 'For compatibility only, not used',
    `name`          VARCHAR(64)          NOT NULL DEFAULT 'unnamed' COMMENT 'The name of the script',
    `source`        TEXT                 NOT NULL                   COMMENT 'Source code for the script',
    `source_ref`    VARCHAR(200)         NOT NULL                   COMMENT 'Where is the script located (if remote)',
    `async`         BOOLEAN              NOT NULL DEFAULT FALSE     COMMENT 'Do we run this script asynchronously?',
    `rel_runner`    BIGINT(20)  UNSIGNED NOT NULL DEFAULT 0         COMMENT 'Who is running the script? 0 for invoker',
    `run_in_ua`     BOOLEAN              NOT NULL DEFAULT FALSE     COMMENT 'Run this script inside user-agent environment',
    `timeout`       INT         UNSIGNED NOT NULL DEFAULT 0         COMMENT 'Any explicit timeout set for this script (milliseconds)?',
    `critical`      BOOLEAN              NOT NULL DEFAULT TRUE      COMMENT 'Is it critical that this script is executed successfully',
    `enabled`       BOOLEAN              NOT NULL DEFAULT TRUE      COMMENT 'Is this script enabled?',

    `created_by`    BIGINT(20)  UNSIGNED NOT NULL DEFAULT 0,
    `created_at`    DATETIME             NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_by`    BIGINT(20)  UNSIGNED NOT NULL DEFAULT 0,
    `updated_at`    DATETIME                 NULL DEFAULT NULL,
    `deleted_by`    BIGINT(20)  UNSIGNED NOT NULL DEFAULT 0,
    `deleted_at`    DATETIME                 NULL DEFAULT NULL,

    PRIMARY KEY (`id`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS sys_automation_trigger (
    `id`         BIGINT(20)  UNSIGNED NOT NULL,
    `rel_script` BIGINT(20)  UNSIGNED NOT NULL              COMMENT 'Script that is triggered',

    `resource`   VARCHAR(128)         NOT NULL              COMMENT 'Resource triggering the event',
    `event`      VARCHAR(128)         NOT NULL              COMMENT 'Event triggered',
    `event_condition`
                 TEXT                 NOT NULL              COMMENT 'Trigger condition',
    `enabled`    BOOLEAN              NOT NULL DEFAULT TRUE COMMENT 'Trigger enabled?',

    `weight`     INT                  NOT NULL DEFAULT 0,

    `created_by` BIGINT(20)  UNSIGNED NOT NULL DEFAULT 0,
    `created_at` DATETIME             NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_by` BIGINT(20)  UNSIGNED NOT NULL DEFAULT 0,
    `updated_at` DATETIME                 NULL DEFAULT NULL,
    `deleted_by` BIGINT(20)  UNSIGNED NOT NULL DEFAULT 0,
    `deleted_at` DATETIME                 NULL DEFAULT NULL,

    CONSTRAINT `fk_sys_automation_script` FOREIGN KEY (`rel_script`) REFERENCES `sys_automation_script` (`id`),

    PRIMARY KEY (`id`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8;
