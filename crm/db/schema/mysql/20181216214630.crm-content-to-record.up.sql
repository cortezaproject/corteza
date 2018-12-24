ALTER TABLE `crm_content` RENAME TO `crm_record`;
ALTER TABLE `crm_record` MODIFY COLUMN `json` json DEFAULT NULL COMMENT 'Records in JSON format.';

ALTER TABLE `crm_content_column` RENAME TO `crm_record_column`;
ALTER TABLE `crm_record_column` CHANGE COLUMN `content_id` `record_id` bigint(20);

ALTER TABLE `crm_content_links` RENAME TO `crm_record_links`;
ALTER TABLE `crm_record_links` CHANGE COLUMN `content_id` `record_id` bigint(20) unsigned;
ALTER TABLE `crm_record_links` CHANGE COLUMN `rel_content_id` `rel_record_id` bigint(20) unsigned;
