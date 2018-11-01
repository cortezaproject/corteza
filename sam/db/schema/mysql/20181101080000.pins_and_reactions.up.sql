DROP TABLE channel_pins;
DROP TABLE reactions;

CREATE TABLE message_flags (
  id               BIGINT UNSIGNED NOT NULL,
  rel_channel      BIGINT UNSIGNED NOT NULL,
  rel_message      BIGINT UNSIGNED NOT NULL,
  rel_user         BIGINT UNSIGNED NOT NULL,
  flag             TEXT,

  created_at       DATETIME        NOT NULL DEFAULT NOW(),

  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
