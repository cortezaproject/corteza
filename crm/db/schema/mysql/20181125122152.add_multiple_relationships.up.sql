CREATE TABLE `crm_content_links` (
 `content_id` bigint(20) unsigned NOT NULL,
 `column_name` varchar(255) NOT NULL,
 `rel_content_id` bigint(20) unsigned NOT NULL,
 PRIMARY KEY (`content_id`,`column_name`,`rel_content_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;