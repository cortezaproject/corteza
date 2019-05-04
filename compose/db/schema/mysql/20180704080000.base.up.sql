CREATE TABLE `crm_content` (
 `id` bigint(20) unsigned NOT NULL,
 `module_id` bigint(20) unsigned NOT NULL,
 `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
 `updated_at` datetime DEFAULT NULL,
 `deleted_at` datetime DEFAULT NULL,
 PRIMARY KEY (`id`,`module_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `crm_content_column` (
 `content_id` bigint(20) NOT NULL,
 `column_name` varchar(255) NOT NULL,
 `column_value` text NOT NULL,
 PRIMARY KEY (`content_id`,`column_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `crm_field` (
 `field_type` varchar(16) NOT NULL COMMENT 'Short field type (string, boolean,...)',
 `field_name` varchar(255) NOT NULL COMMENT 'Description of field contents',
 `field_template` varchar(255) NOT NULL COMMENT 'HTML template file for field',
 PRIMARY KEY (`field_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `crm_module` (
 `id` bigint(20) unsigned NOT NULL,
 `name` varchar(64) NOT NULL COMMENT 'The name of the module',
 `json` json NOT NULL COMMENT 'List of field definitions for the module',
 `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
 `updated_at` datetime DEFAULT NULL,
 `deleted_at` datetime DEFAULT NULL,
 PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `crm_module_form` (
 `module_id` bigint(20) unsigned NOT NULL,
 `place` tinyint(3) unsigned NOT NULL,
 `kind` varchar(64) NOT NULL COMMENT 'The type of the form input field',
 `name` varchar(64) NOT NULL COMMENT 'The name of the field in the form',
 `label` varchar(255) NOT NULL COMMENT 'The label of the form input',
 `help_text` text NOT NULL COMMENT 'Help text',
 `default_value` text NOT NULL COMMENT 'Default value',
 `max_length` int(10) unsigned NOT NULL COMMENT 'Maximum input length',
 `is_private` tinyint(1) NOT NULL COMMENT 'Contains personal/sensitive data?',
 PRIMARY KEY (`module_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `crm_page` (
 `id` bigint(20) unsigned NOT NULL COMMENT 'Page ID',
 `self_id` bigint(20) unsigned NOT NULL COMMENT 'Parent Page ID',
 `module_id` bigint(20) unsigned NOT NULL COMMENT 'Module ID (optional)',
 `title` varchar(255) NOT NULL COMMENT 'Title (required)',
 `description` text NOT NULL COMMENT 'Description',
 `blocks` json NOT NULL COMMENT 'JSON array of blocks for the page',
 `visible` tinyint(4) NOT NULL COMMENT 'Is page visible in navigation?',
 `weight` int(11) NOT NULL COMMENT 'Order for navigation',
 PRIMARY KEY (`id`) USING BTREE,
 KEY `module_id` (`module_id`),
 KEY `self_id` (`self_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

