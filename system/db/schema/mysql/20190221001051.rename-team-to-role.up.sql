ALTER TABLE sys_team RENAME TO sys_role;
ALTER TABLE sys_team_member RENAME TO sys_role_member;

ALTER TABLE `sys_role_member` CHANGE COLUMN `rel_team` `rel_role` BIGINT UNSIGNED NOT NULL;
ALTER TABLE `sys_rules` CHANGE COLUMN `rel_team` `rel_role` BIGINT UNSIGNED NOT NULL;
