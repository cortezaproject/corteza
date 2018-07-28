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
);

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
);

-- Keeps all known channels
CREATE TABLE channels (
  id               BIGINT UNSIGNED NOT NULL,
  name             TEXT            NOT NULL, -- display name of the channel
  meta             JSON            NOT NULL,

  created_at       DATETIME        NOT NULL DEFAULT NOW(),
  updated_at       DATETIME            NULL,
  archived_at      DATETIME            NULL,
  deleted_at       DATETIME            NULL, -- channel soft delete

  rel_last_message BIGINT UNSIGNED,

  PRIMARY KEY (id)
);

-- Keeps all known users, home and external organisation
--   changes are stored in audit log
CREATE TABLE users (
  id               BIGINT UNSIGNED NOT NULL,
  username         TEXT            NOT NULL,
  password         TEXT,
  meta             JSON            NOT NULL,

  rel_organisation BIGINT UNSIGNED NOT NULL REFERENCES organisation(id),

  created_at       DATETIME        NOT NULL DEFAULT NOW(),
  updated_at       DATETIME            NULL,
  suspended_at     DATETIME            NULL,
  deleted_at       DATETIME            NULL, -- user soft delete

  PRIMARY KEY (id)
);

-- Keeps team memberships
CREATE TABLE team_members (
  rel_team         BIGINT UNSIGNED NOT NULL REFERENCES organisation(id),
  rel_user         BIGINT UNSIGNED NOT NULL REFERENCES users(id),

  PRIMARY KEY (rel_team, rel_user)
);

-- handles channel membership
CREATE TABLE channel_members (
  rel_channel      BIGINT UNSIGNED NOT NULL REFERENCES channels(id),
  rel_user         BIGINT UNSIGNED NOT NULL REFERENCES users(id),

  PRIMARY KEY (rel_channel, rel_user)
);

CREATE TABLE channel_views (
  rel_channel      BIGINT UNSIGNED NOT NULL REFERENCES channels(id),
  rel_user         BIGINT UNSIGNED NOT NULL REFERENCES users(id),

  -- timestamp of last view, should be enough to find out which messaghr
  viewed_at        DATETIME        NOT NULL DEFAULT NOW(),

  -- new messages count since last view
  new_since        INT    UNSIGNED NOT NULL DEFAULT 0,

  PRIMARY KEY (rel_user, rel_channel)
);

CREATE TABLE channel_pins (
  rel_channel      BIGINT UNSIGNED NOT NULL REFERENCES channels(id),
  rel_message      BIGINT UNSIGNED NOT NULL REFERENCES messages(id),
  rel_user         BIGINT UNSIGNED NOT NULL REFERENCES users(id),

  created_at       DATETIME        NOT NULL DEFAULT NOW(),

  PRIMARY KEY (rel_channel, rel_message)
);

CREATE TABLE messages (
  id               BIGINT UNSIGNED NOT NULL,
  type             TEXT,
  message          TEXT            NOT NULL,
  meta             JSON,
  rel_user         BIGINT UNSIGNED NOT NULL REFERENCES users(id),
  rel_channel      BIGINT UNSIGNED NOT NULL REFERENCES channels(id),
  reply_to         BIGINT UNSIGNED     NULL REFERENCES messages(id),

  created_at       DATETIME        NOT NULL DEFAULT NOW(),
  updated_at       DATETIME            NULL,
  deleted_at       DATETIME            NULL,

  PRIMARY KEY (id)
);

CREATE TABLE reactions (
  id               BIGINT UNSIGNED NOT NULL,
  rel_user         BIGINT UNSIGNED NOT NULL REFERENCES users(id),
  rel_message      BIGINT UNSIGNED NOT NULL REFERENCES messages(id),
  rel_channel      BIGINT UNSIGNED NOT NULL REFERENCES channels(id),
  reaction         TEXT            NOT NULL,

  created_at       DATETIME        NOT NULL DEFAULT NOW(),

  PRIMARY KEY (id)
);

CREATE TABLE attachments (
  id               BIGINT UNSIGNED NOT NULL,
  rel_user         BIGINT UNSIGNED NOT NULL REFERENCES users(id),
  rel_message      BIGINT UNSIGNED NOT NULL REFERENCES messages(id),
  rel_channel      BIGINT UNSIGNED NOT NULL REFERENCES channels(id),
  url              TEXT,
  preview_url      TEXT,
  size             INT    UNSIGNED,
  mimetype         TEXT,
  name             TEXT,
  attachment       JSON            NOT NULL,

  created_at       DATETIME        NOT NULL DEFAULT NOW(),
  updated_at       DATETIME            NULL,
  deleted_at       DATETIME            NULL,

  PRIMARY KEY (id)
);
