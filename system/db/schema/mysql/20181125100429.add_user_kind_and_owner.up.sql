# add field to manage user type (bot support)
ALTER TABLE `sys_user` ADD `kind` VARCHAR(8) NOT NULL DEFAULT '' AFTER `handle`;

# add field to manage "ownership" (get all bots created by user)
ALTER TABLE `sys_user` ADD `rel_user_id` BIGINT UNSIGNED NOT NULL AFTER `rel_organisation`, ADD INDEX (`rel_user_id`);
