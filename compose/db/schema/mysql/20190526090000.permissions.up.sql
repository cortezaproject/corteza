CREATE TABLE IF NOT EXISTS compose_permission_rules (
  rel_role   BIGINT UNSIGNED NOT NULL,
  resource   VARCHAR(128)    NOT NULL,
  operation  VARCHAR(128)    NOT NULL,
  access     TINYINT(1)      NOT NULL,

  PRIMARY KEY (rel_role, resource, operation)
) ENGINE=InnoDB;
