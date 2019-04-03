ALTER TABLE `sys_user` DROP `password`;
ALTER TABLE `sys_user` DROP `satosa_id`;
ALTER TABLE `sys_credentials` ADD `last_used_at` DATETIME NULL;
