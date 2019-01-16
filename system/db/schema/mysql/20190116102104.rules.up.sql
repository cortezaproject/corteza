CREATE TABLE `sys_rules` (
  `rel_team` BIGINT UNSIGNED NOT NULL,
  `resource` VARCHAR(128) NOT NULL,
  `operation` VARCHAR(128) NOT NULL,
  `value` ENUM('no', 'yes') NOT NULL,

  PRIMARY KEY (`rel_team`, `resource`, `operation`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
