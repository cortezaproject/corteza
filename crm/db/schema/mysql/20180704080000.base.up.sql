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

CREATE TABLE `crm_module_content` (
 `id` bigint(20) unsigned NOT NULL,
 `module_id` bigint(20) unsigned NOT NULL,
 `json` json NOT NULL COMMENT 'Entry JSON blob based on module fields list',
 PRIMARY KEY (`id`,`module_id`)
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

INSERT INTO `crm_fields` VALUES ('bool','Boolean value (yes / no)','');
INSERT INTO `crm_fields` VALUES ('email','E-mail input','');
INSERT INTO `crm_fields` VALUES ('enum','Single option picker','');
INSERT INTO `crm_fields` VALUES ('hidden','Hidden value','');
INSERT INTO `crm_fields` VALUES ('stamp','Date/time input','');
INSERT INTO `crm_fields` VALUES ('text','Text input','');
INSERT INTO `crm_fields` VALUES ('textarea','Text input (multi-line)','');
