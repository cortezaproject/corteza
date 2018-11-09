ALTER TABLE `crm_content` ADD `user_id` BIGINT UNSIGNED NOT NULL AFTER `module_id`, ADD INDEX (`user_id`);
