DROP TABLE IF EXISTS compose_record;
CREATE TABLE `compose_record` (
  `id` bigint unsigned NOT NULL,
  `rel_namespace` bigint unsigned NOT NULL,
  `module_id` bigint unsigned NOT NULL,

  `values` json DEFAULT NULL,

  `owned_by` bigint unsigned NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `created_by` bigint unsigned NOT NULL,
  `updated_by` bigint unsigned NOT NULL DEFAULT '0',
  `deleted_by` bigint unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `compose_record_owner` (`owned_by`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;

DROP TABLE IF EXISTS compose_record_partitioned;
CREATE TABLE `compose_record_partitioned` (
  `id` bigint unsigned NOT NULL,

  `vID` BIGINT UNSIGNED,
  `vRef` BIGINT UNSIGNED,
  `vTimestamp` DATETIME,
  `vTime` TIME,
  `vDate` DATE,
  `vNumber` NUMERIC,
  `vText` TEXT,
  `vBoolean_T` BOOLEAN,
  `vBoolean_F` BOOLEAN,
  `vEnum` TEXT,
  `vGeometry` TEXT,
  `vJSON` TEXT,
  `vBlob` BLOB,
  `vUUID` VARCHAR(36),

  `owned_by` bigint unsigned NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `created_by` bigint unsigned NOT NULL,
  `updated_by` bigint unsigned NOT NULL DEFAULT '0',
  `deleted_by` bigint unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `compose_record_owner` (`owned_by`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;
