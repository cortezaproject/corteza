-- No more links, we'll handle this through ref field on crm_record_value tbl
DROP TABLE IF EXISTS `crm_record_links`;

-- Not columns, values
ALTER TABLE `crm_record_column` RENAME TO `crm_record_value`;

-- Simplify names
ALTER TABLE `crm_record_value` CHANGE COLUMN `column_name`  `name`  VARCHAR(64);
ALTER TABLE `crm_record_value` CHANGE COLUMN `column_value` `value` TEXT;

-- Add reference
ALTER TABLE `crm_record_value` ADD  COLUMN `ref` BIGINT UNSIGNED DEFAULT 0 NOT NULL;
ALTER TABLE `crm_record_value` ADD  COLUMN `deleted_at` datetime DEFAULT NULL;
ALTER TABLE `crm_record_value` ADD  COLUMN `place` INT UNSIGNED DEFAULT 0 NOT NULL;
ALTER TABLE `crm_record_value` DROP PRIMARY KEY, ADD PRIMARY KEY(`record_id`, `name`, `place`);
CREATE INDEX crm_record_value_ref ON crm_record_value (ref);


-- We want this as a real field
ALTER TABLE `crm_module_form`  ADD  COLUMN `is_multi` TINYINT(1) NOT NULL;

-- This will be handled through meta(json) fieldd
ALTER TABLE `crm_module_form`  DROP COLUMN `help_text`;
ALTER TABLE `crm_module_form`  DROP COLUMN `max_length`;
ALTER TABLE `crm_module_form`  DROP COLUMN `default_Value`;
