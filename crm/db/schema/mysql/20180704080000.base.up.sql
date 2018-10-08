CREATE TABLE `crm_content` (
 `id` bigint(20) unsigned NOT NULL,
 `module_id` bigint(20) unsigned NOT NULL,
 `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
 `updated_at` datetime DEFAULT NULL,
 `archived_at` datetime DEFAULT NULL,
 `deleted_at` datetime DEFAULT NULL,
 PRIMARY KEY (`id`,`module_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `crm_content_columns` (
 `content_id` bigint(20) NOT NULL,
 `column_name` varchar(255) NOT NULL,
 `column_value` text NOT NULL,
 PRIMARY KEY (`content_id`,`column_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `crm_fields` (
 `field_type` varchar(16) NOT NULL COMMENT 'Short field type (string, boolean,...)',
 `field_name` varchar(255) NOT NULL COMMENT 'Description of field contents',
 `field_template` varchar(255) NOT NULL COMMENT 'HTML template file for field',
 PRIMARY KEY (`field_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `crm_module` (
 `id` bigint(20) unsigned NOT NULL,
 `name` varchar(64) NOT NULL COMMENT 'The name of the module',
 `json` json NOT NULL COMMENT 'List of field definitions for the module',
 PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `crm_module_form` (
 `module_id` bigint(20) unsigned NOT NULL,
 `place` tinyint(3) unsigned NOT NULL,
 `type` varchar(64) NOT NULL COMMENT 'The type of the form input field',
 `name` varchar(64) NOT NULL COMMENT 'The name of the field in the form',
 `title` varchar(255) NOT NULL COMMENT 'The label of the form input',
 `placeholder` varchar(255) NOT NULL COMMENT 'Placeholder text (if supported by input)',
 PRIMARY KEY (`module_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

