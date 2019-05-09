DROP TABLE IF EXISTS sys_permission_rules;
DROP TABLE IF EXISTS messaging_permission_rules;
DROP TABLE IF EXISTS compose_permission_rules;

CREATE TABLE IF NOT EXISTS sys_permission_rules (
  rel_role   BIGINT UNSIGNED NOT NULL,
  resource   VARCHAR(128)    NOT NULL,
  operation  VARCHAR(128)    NOT NULL,
  access     TINYINT(1)      NOT NULL,

  PRIMARY KEY (rel_role, resource, operation)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS messaging_permission_rules (
  rel_role   BIGINT UNSIGNED NOT NULL,
  resource   VARCHAR(128)    NOT NULL,
  operation  VARCHAR(128)    NOT NULL,
  access     TINYINT(1)      NOT NULL,

  PRIMARY KEY (rel_role, resource, operation)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS compose_permission_rules (
  rel_role   BIGINT UNSIGNED NOT NULL,
  resource   VARCHAR(128)    NOT NULL,
  operation  VARCHAR(128)    NOT NULL,
  access     TINYINT(1)      NOT NULL,

  PRIMARY KEY (rel_role, resource, operation)
) ENGINE=InnoDB;

REPLACE sys_permission_rules
    (rel_role, resource, operation, access)
    SELECT rel_role, resource, operation, `value` - 1 FROM sys_rules WHERE resource LIKE 'system%';

REPLACE compose_permission_rules
    (rel_role, resource, operation, access)
    SELECT rel_role, resource, operation, `value` - 1 FROM sys_rules WHERE resource LIKE 'compose%';

REPLACE messaging_permission_rules
    (rel_role, resource, operation, access)
    SELECT rel_role, resource, operation, `value` - 1 FROM sys_rules WHERE resource LIKE 'messaging%';

DROP TABLE sys_rules;
