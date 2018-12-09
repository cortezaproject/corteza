-- Keeps all known users, home and external organisation
--   changes are stored in audit log
CREATE TABLE sys_credentials (
  id               BIGINT UNSIGNED NOT NULL,
  rel_owner        BIGINT UNSIGNED NOT NULL REFERENCES sys_users(id),
  label            TEXT            NOT NULL COMMENT 'something we can differentiate credentials by',
  kind             VARCHAR(128)    NOT NULL COMMENT 'hash, facebook, gplus, github, linkedin ...',
  credentials      TEXT            NOT NULL COMMENT 'crypted/hashed passwords, secrets, social profile ID',
  meta             JSON            NOT NULL,
  expires_at       DATETIME            NULL,

  created_at       DATETIME        NOT NULL DEFAULT NOW(),
  updated_at       DATETIME            NULL,
  deleted_at       DATETIME            NULL, -- user soft delete

  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE INDEX idx_owner ON sys_credentials (rel_owner);
