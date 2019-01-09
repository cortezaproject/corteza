CREATE TABLE `crm_trigger` (
 `id`         BIGINT(20)  UNSIGNED NOT NULL,
 `name`       VARCHAR(64)          NOT NULL COMMENT 'The name of the trigger',
 `enabled`    BOOLEAN              NOT NULL COMMENT 'Trigger enabled?',
 `actions`    TEXT                 NOT NULL COMMENT 'All actions that trigger it',
 `source`     TEXT                 NOT NULL COMMENT 'Trigger source',
 `rel_module` BIGINT(20)  UNSIGNED     NULL COMMENT 'Primary module',

 `created_at` DATETIME             NOT NULL DEFAULT CURRENT_TIMESTAMP,
 `updated_at` DATETIME                      DEFAULT NULL,
 `deleted_at` DATETIME                      DEFAULT NULL,

 PRIMARY KEY (`id`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8;
