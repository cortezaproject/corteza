CREATE TABLE `compose_namespace` (
 `id`         BIGINT(20)  UNSIGNED NOT NULL,
 `name`       VARCHAR(64)          NOT NULL COMMENT 'Name',
 `slug`       VARCHAR(64)          NOT NULL COMMENT 'URL slug',
 `enabled`    BOOLEAN              NOT NULL COMMENT 'Is namespace enabled?',
 `meta`       JSON                 NOT NULL COMMENT 'Meta data',

 `created_at` DATETIME             NOT NULL DEFAULT CURRENT_TIMESTAMP,
 `updated_at` DATETIME                      DEFAULT NULL,
 `deleted_at` DATETIME                      DEFAULT NULL,

 PRIMARY KEY (`id`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8;
