-- Keeps all known channels
CREATE TABLE channels (
  id               BIGINT UNSIGNED NOT NULL,
  name             TEXT            NOT NULL, -- display name of the channel
  topic            TEXT            NOT NULL,
  meta             JSON            NOT NULL,

  type             ENUM ('private', 'public', 'group') NOT NULL DEFAULT 'public',

  rel_organisation BIGINT UNSIGNED NOT NULL REFERENCES organisation(id),
  rel_creator      BIGINT UNSIGNED NOT NULL,

  created_at       DATETIME        NOT NULL DEFAULT NOW(),
  updated_at       DATETIME            NULL,
  archived_at      DATETIME            NULL,
  deleted_at       DATETIME            NULL, -- channel soft delete

  rel_last_message BIGINT UNSIGNED NOT NULL DEFAULT 0,

  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- handles channel membership
CREATE TABLE channel_members (
  rel_channel      BIGINT UNSIGNED NOT NULL REFERENCES channels(id),
  rel_user         BIGINT UNSIGNED NOT NULL,

  type             ENUM ('owner', 'member', 'invitee') NOT NULL DEFAULT 'member',

  created_at       DATETIME        NOT NULL DEFAULT NOW(),
  updated_at       DATETIME            NULL,

  PRIMARY KEY (rel_channel, rel_user)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE channel_views (
  rel_channel      BIGINT UNSIGNED NOT NULL REFERENCES channels(id),
  rel_user         BIGINT UNSIGNED NOT NULL,

  -- timestamp of last view, should be enough to find out which messaghr
  viewed_at        DATETIME        NOT NULL DEFAULT NOW(),

  -- new messages count since last view
  new_since        INT    UNSIGNED NOT NULL DEFAULT 0,

  PRIMARY KEY (rel_user, rel_channel)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE channel_pins (
  rel_channel      BIGINT UNSIGNED NOT NULL REFERENCES channels(id),
  rel_message      BIGINT UNSIGNED NOT NULL REFERENCES messages(id),
  rel_user         BIGINT UNSIGNED NOT NULL,

  created_at       DATETIME        NOT NULL DEFAULT NOW(),

  PRIMARY KEY (rel_channel, rel_message)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE messages (
  id               BIGINT UNSIGNED NOT NULL,
  type             TEXT,
  message          TEXT            NOT NULL,
  meta             JSON,
  rel_user         BIGINT UNSIGNED NOT NULL,
  rel_channel      BIGINT UNSIGNED NOT NULL REFERENCES channels(id),
  reply_to         BIGINT UNSIGNED     NULL REFERENCES messages(id),

  created_at       DATETIME        NOT NULL DEFAULT NOW(),
  updated_at       DATETIME            NULL,
  deleted_at       DATETIME            NULL,

  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE reactions (
  id               BIGINT UNSIGNED NOT NULL,
  rel_user         BIGINT UNSIGNED NOT NULL,
  rel_message      BIGINT UNSIGNED NOT NULL REFERENCES messages(id),
  rel_channel      BIGINT UNSIGNED NOT NULL REFERENCES channels(id),
  reaction         TEXT            NOT NULL,

  created_at       DATETIME        NOT NULL DEFAULT NOW(),

  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE attachments (
  id               BIGINT UNSIGNED NOT NULL,
  rel_user         BIGINT UNSIGNED NOT NULL,

  url              VARCHAR(512),
  preview_url      VARCHAR(512),

  size             INT    UNSIGNED,
  mimetype         VARCHAR(255),
  name             TEXT,

  meta             JSON,

  created_at       DATETIME        NOT NULL DEFAULT NOW(),
  updated_at       DATETIME            NULL,
  deleted_at       DATETIME            NULL,

  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE message_attachment (
  rel_message      BIGINT UNSIGNED NOT NULL REFERENCES messages(id),
  rel_attachment   BIGINT UNSIGNED NOT NULL REFERENCES attachment(id),

  PRIMARY KEY (rel_message)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE event_queue (
  id               BIGINT UNSIGNED NOT NULL,
  origin           BIGINT UNSIGNED NOT NULL,
  subscriber       TEXT,
  payload          JSON,

  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE event_queue_synced (
  origin           BIGINT UNSIGNED NOT NULL,
  rel_last         BIGINT UNSIGNED NOT NULL,

  PRIMARY KEY (origin)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
