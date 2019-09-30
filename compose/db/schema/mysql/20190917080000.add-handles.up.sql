ALTER TABLE `compose_module` ADD `handle` VARCHAR(200) NOT NULL AFTER `id`;
ALTER TABLE `compose_page`   ADD `handle` VARCHAR(200) NOT NULL AFTER `id`;
ALTER TABLE `compose_chart`  ADD `handle` VARCHAR(200) NOT NULL AFTER `id`;
