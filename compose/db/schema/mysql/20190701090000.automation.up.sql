DROP TABLE IF EXISTS compose_automation_trigger;
DROP TABLE IF EXISTS compose_automation_script;

CREATE TABLE IF NOT EXISTS compose_automation_script (
    `id`         BIGINT(20)  UNSIGNED NOT NULL,
    `name`       VARCHAR(64)          NOT NULL DEFAULT 'unnamed' COMMENT 'The name of the script',
    `source`     TEXT                 NOT NULL                   COMMENT 'Source code for the script',
    `source_ref` VARCHAR(200)         NOT NULL                   COMMENT 'Where is the script located (if remote)',
    `async`      BOOLEAN              NOT NULL DEFAULT FALSE     COMMENT 'Do we run this script asynchronously?',
    `rel_runner` BIGINT(20)  UNSIGNED NOT NULL DEFAULT 0         COMMENT 'Who is running the script? 0 for invoker',
    `run_in_ua`  BOOLEAN              NOT NULL DEFAULT FALSE     COMMENT 'Run this script inside user-agent environment',
    `timeout`    INT         UNSIGNED NOT NULL DEFAULT 0         COMMENT 'Any explicit timeout set for this script (milliseconds)?',
    `critical`   BOOLEAN              NOT NULL DEFAULT TRUE      COMMENT 'Is it critical that this script is executed successfully',
    `enabled`    BOOLEAN              NOT NULL DEFAULT TRUE      COMMENT 'Is this script enabled?',

    `created_by` BIGINT(20)  UNSIGNED NOT NULL DEFAULT 0,
    `created_at` DATETIME             NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_by` BIGINT(20)  UNSIGNED NOT NULL DEFAULT 0,
    `updated_at` DATETIME                 NULL DEFAULT NULL,
    `deleted_by` BIGINT(20)  UNSIGNED NOT NULL DEFAULT 0,
    `deleted_at` DATETIME                 NULL DEFAULT NULL,

    PRIMARY KEY (`id`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS compose_automation_trigger (
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

    CONSTRAINT `fk_script` FOREIGN KEY (`rel_script`) REFERENCES `compose_automation_script` (`id`),

    PRIMARY KEY (`id`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Migrate old triggers into scripts
INSERT INTO compose_automation_script (id, name, source, source_ref, run_in_ua, critical, enabled, created_at, updated_at, deleted_at)
SELECT id, name, source, '', true, false, enabled, created_at, updated_at, deleted_at from compose_trigger;

# Migrate old triggers into new triggers
INSERT INTO compose_automation_trigger (id, event, resource, event_condition, rel_script, enabled, created_at, updated_at, deleted_at)
SELECT id+seq, events.event, 'compose:record', rel_module, id, enabled, created_at, updated_at, deleted_at from compose_trigger AS t INNER JOIN
              (      SELECT 0 as seq, ''             AS event
               UNION SELECT 1 as seq, 'manual'       AS event
               UNION SELECT 2 as seq, 'beforeCreate' AS event
               UNION SELECT 3 as seq, 'afterCreate'  AS event
               UNION SELECT 4 as seq, 'beforeUpdate' AS event
               UNION SELECT 5 as seq, 'afterUpdate'  AS event
               UNION SELECT 6 as seq, 'beforeDelete' AS event
               UNION SELECT 7 as seq, 'afterDelete'  AS event) AS events ON ((event  = '' AND t.actions = '')
                                                                          OR (event <> '' AND t.actions LIKE concat('%',event,'%') ));
# Normalize and cleanup
UPDATE compose_automation_trigger SET event = 'manual' WHERE event = '';
DELETE FROM compose_automation_trigger WHERE event_condition IN ('', '0') AND event <> 'manual';

DROP TABLE IF EXISTS compose_trigger;
