ALTER TABLE `messaging_channel` ADD `membership_policy` ENUM ('featured', 'forced', '') NOT NULL DEFAULT '' AFTER `type`;
