-- all known organisations (crust instances) and our relation towards them
CREATE TABLE organisations (
  id               BIGINT UNSIGNED NOT NULL,
  fqn              TEXT            NOT NULL, -- fully qualified name of the organisation
  name             TEXT            NOT NULL, -- display name of the organisation

  created_at       DATETIME        NOT NULL DEFAULT NOW(),
  updated_at       DATETIME            NULL,
  archived_at      DATETIME            NULL,
  deleted_at       DATETIME            NULL, -- organisation soft delete

  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

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
  satosa_id        CHAR(36)            NULL,

  rel_organisation BIGINT UNSIGNED NOT NULL,

  created_at       DATETIME        NOT NULL DEFAULT NOW(),
  updated_at       DATETIME            NULL,
  suspended_at     DATETIME            NULL,
  deleted_at       DATETIME            NULL, -- user soft delete

  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE UNIQUE INDEX uid_satosa ON users (satosa_id);

-- Keeps all known teams
CREATE TABLE teams (
  id               BIGINT UNSIGNED NOT NULL,
  name             TEXT            NOT NULL, -- display name of the team
  handle           TEXT            NOT NULL, -- team handle string

  created_at       DATETIME        NOT NULL DEFAULT NOW(),
  updated_at       DATETIME            NULL,
  archived_at      DATETIME            NULL,
  deleted_at       DATETIME            NULL, -- team soft delete

  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Keeps team memberships
CREATE TABLE team_members (
  rel_team         BIGINT UNSIGNED NOT NULL REFERENCES organisation(id),
  rel_user         BIGINT UNSIGNED NOT NULL,

  PRIMARY KEY (rel_team, rel_user)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
