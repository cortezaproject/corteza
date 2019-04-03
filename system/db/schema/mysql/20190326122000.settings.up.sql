DROP TABLE `settings`;

CREATE TABLE `sys_settings` (
  rel_owner        BIGINT UNSIGNED NOT NULL DEFAULT 0     COMMENT 'Value owner, 0 for global settings',
  name             VARCHAR(200)    NOT NULL               COMMENT 'Unique set of setting keys',
  value            JSON                                   COMMENT 'Setting value',

  updated_at       DATETIME        NOT NULL DEFAULT NOW() COMMENT 'When was the value updated',
  updated_by       BIGINT UNSIGNED NOT NULL DEFAULT 0     COMMENT 'Who created/updated the value',

  PRIMARY KEY (name, rel_owner)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
