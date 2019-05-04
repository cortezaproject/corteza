ALTER TABLE `crm_record` CHANGE COLUMN `user_id`  `owned_by` BIGINT UNSIGNED NOT NULL DEFAULT 0;
ALTER TABLE `crm_record` ADD COLUMN `created_by` BIGINT UNSIGNED NOT NULL DEFAULT 0;
ALTER TABLE `crm_record` ADD COLUMN `updated_by` BIGINT UNSIGNED NOT NULL DEFAULT 0;
ALTER TABLE `crm_record` ADD COLUMN `deleted_by` BIGINT UNSIGNED NOT NULL DEFAULT 0;
UPDATE crm_record SET created_by = owned_by;
UPDATE crm_record SET updated_by = owned_by WHERE updated_at IS NOT NULL;
UPDATE crm_record SET deleted_by = owned_by WHERE deleted_at IS NOT NULL;
