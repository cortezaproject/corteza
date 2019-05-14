ALTER TABLE compose_module_form
    RENAME TO compose_module_field;

-- Remove orphaned and invalid fields
DELETE FROM `compose_module_field` WHERE `module_id` NOT IN (SELECT `id` FROM `compose_module`) OR `name` = '';

-- Order and consistency.
ALTER TABLE `compose_module_field`
    ADD COLUMN `id`         BIGINT UNSIGNED NOT NULL FIRST,
    ADD COLUMN `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ADD COLUMN `updated_at` DATETIME DEFAULT NULL,
    ADD COLUMN `deleted_at` DATETIME DEFAULT NULL,
    RENAME COLUMN `module_id` TO `rel_module`,
    RENAME COLUMN `json`      TO `options`;

-- Generate IDs for the new field, use module, offset by one (just to start with a different ID)
-- and use place (0 based, +1 for every field, expecting to be unique per module because of the existing pkey)
UPDATE `compose_module_field` SET id = rel_module + 1 + place;

-- Drop old primary key (module_id, place)
ALTER TABLE `compose_module_field` DROP PRIMARY KEY, ADD PRIMARY KEY(`id`);

-- Foreign key
ALTER TABLE `compose_module_field`
    ADD CONSTRAINT `compose_module`
        FOREIGN KEY (`rel_module`)
            REFERENCES `compose_module` (`id`);

-- And unique indexes for module+place/name combos.
CREATE UNIQUE INDEX uid_compose_module_field_place ON compose_module_field (`rel_module`, `place`);
CREATE UNIQUE INDEX uid_compose_module_field_name  ON compose_module_field (`rel_module`, `name`);
