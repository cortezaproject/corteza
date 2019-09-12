ALTER TABLE `compose_module_field`
  ADD `default_value` JSON DEFAULT NULL COMMENT 'Default value as a record value set.'
  AFTER `options`;
