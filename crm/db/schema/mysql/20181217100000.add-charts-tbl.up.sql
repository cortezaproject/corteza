CREATE TABLE `crm_chart` (
 `id`         BIGINT(20)  UNSIGNED NOT NULL,
 `name`       VARCHAR(64)          NOT NULL COMMENT 'The name of the chart',
 `config`     JSON                 NOT NULL COMMENT 'Chart & reporting configuration',

 `created_at` DATETIME             NOT NULL DEFAULT CURRENT_TIMESTAMP,
 `updated_at` DATETIME                      DEFAULT NULL,
 `deleted_at` DATETIME                      DEFAULT NULL,

 PRIMARY KEY (`id`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8;
