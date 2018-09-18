CREATE TABLE settings (
  name  VARCHAR(200) NOT NULL   COMMENT 'Unique set of setting keys',
  value TEXT                    COMMENT 'Setting value',

  PRIMARY KEY (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Keeps all known users, home and external organisation
--   changes are stored in audit log
CREATE TABLE users (
  id               BIGINT UNSIGNED NOT NULL,
  email            TEXT            NOT NULL,
  username         TEXT            NOT NULL,
  password         TEXT            NOT NULL,
  name             TEXT            NOT NULL,
  handle           TEXT            NOT NULL,
  meta             JSON            NOT NULL,

  rel_organisation BIGINT UNSIGNED NOT NULL,

  created_at       DATETIME        NOT NULL DEFAULT NOW(),
  updated_at       DATETIME            NULL,
  suspended_at     DATETIME            NULL,
  deleted_at       DATETIME            NULL, -- user soft delete

  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
