ALTER TABLE `compose_automation_script`
    ADD `rel_namespace` BIGINT UNSIGNED NOT NULL AFTER `id`,
    ADD INDEX (`rel_namespace`);

UPDATE `compose_automation_script` SET `rel_namespace` = (SELECT MIN(id) FROM compose_namespace);

ALTER TABLE `compose_automation_script`
    ADD CONSTRAINT `compose_automation_script_namespace`
    FOREIGN KEY (`rel_namespace`)
    REFERENCES `compose_namespace` (`id`);
