CREATE TABLE `messaging_webhook` (
 `id` bigint(20) unsigned NOT NULL,
 `kind` varchar(8) NOT NULL COMMENT 'Kind: incoming, outgoing',
 `token` varchar(255) NOT NULL COMMENT 'Authentication token',
 `rel_owner` bigint(20) unsigned NOT NULL COMMENT 'Webhook owner User ID',
 `rel_user` bigint(20) unsigned NOT NULL COMMENT 'Webhook message User ID',
 `rel_channel` bigint(20) unsigned NOT NULL COMMENT 'Channel ID',
 `outgoing_trigger` varchar(32) NOT NULL COMMENT 'Outgoing command trigger',
 `outgoing_url` varchar(255) NOT NULL COMMENT 'URL for POST request',
 `created_at` datetime NOT NULL,
 `updated_at` datetime     NULL,
 `deleted_at` datetime     NULL,
 PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- get webhook by command trigger
ALTER TABLE `messaging_webhook` ADD UNIQUE(`outgoing_trigger`);

-- list webhooks by owner (list your own webhooks)
ALTER TABLE `messaging_webhook` ADD INDEX(`rel_owner`);

-- list webhooks on a channel
ALTER TABLE `messaging_webhook` ADD INDEX(`rel_channel`);
