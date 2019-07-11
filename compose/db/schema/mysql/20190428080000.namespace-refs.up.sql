INSERT IGNORE INTO compose_namespace (id, name, slug, enabled, meta)
            VALUES (88714882739863655, 'Crust CRM', 'crm', true, '{}');

ALTER TABLE `compose_attachment`
        ADD `rel_namespace` BIGINT UNSIGNED NOT NULL AFTER `id`,
        ADD INDEX (`rel_namespace`);

ALTER TABLE `compose_chart`
        ADD `rel_namespace` BIGINT UNSIGNED NOT NULL AFTER `id`,
        ADD INDEX (`rel_namespace`);

ALTER TABLE `compose_module`
        ADD `rel_namespace` BIGINT UNSIGNED NOT NULL AFTER `id`,
        ADD INDEX (`rel_namespace`);

ALTER TABLE `compose_page`
        ADD `rel_namespace` BIGINT UNSIGNED NOT NULL AFTER `id`,
        ADD INDEX (`rel_namespace`);

ALTER TABLE `compose_record`
        ADD `rel_namespace` BIGINT UNSIGNED NOT NULL AFTER `id`,
        ADD INDEX (`rel_namespace`);

ALTER TABLE `compose_trigger`
        ADD `rel_namespace` BIGINT UNSIGNED NOT NULL AFTER `id`,
        ADD INDEX (`rel_namespace`);

UPDATE `compose_attachment`   SET `rel_namespace` = 88714882739863655;
UPDATE `compose_chart`        SET `rel_namespace` = 88714882739863655;
UPDATE `compose_module`       SET `rel_namespace` = 88714882739863655;
UPDATE `compose_page`         SET `rel_namespace` = 88714882739863655;
UPDATE `compose_record`       SET `rel_namespace` = 88714882739863655;
UPDATE `compose_trigger`      SET `rel_namespace` = 88714882739863655;


ALTER TABLE `compose_attachment`
        ADD CONSTRAINT `compose_attachment_namespace`
            FOREIGN KEY (`rel_namespace`)
            REFERENCES `compose_namespace` (`id`);

ALTER TABLE `compose_chart`
        ADD CONSTRAINT `compose_chart_namespace`
            FOREIGN KEY (`rel_namespace`)
            REFERENCES `compose_namespace` (`id`);

ALTER TABLE `compose_module`
        ADD CONSTRAINT `compose_module_namespace`
            FOREIGN KEY (`rel_namespace`)
            REFERENCES `compose_namespace` (`id`);

ALTER TABLE `compose_page`
        ADD CONSTRAINT `compose_page_namespace`
            FOREIGN KEY (`rel_namespace`)
            REFERENCES `compose_namespace` (`id`);

ALTER TABLE `compose_record`
        ADD CONSTRAINT `compose_record_namespace`
            FOREIGN KEY (`rel_namespace`)
            REFERENCES `compose_namespace` (`id`);

ALTER TABLE `compose_trigger`
        ADD CONSTRAINT `compose_trigger_namespace`
            FOREIGN KEY (`rel_namespace`)
            REFERENCES `compose_namespace` (`id`);
