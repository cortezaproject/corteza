CREATE TABLE mentions (
  id               BIGINT UNSIGNED NOT NULL,
  rel_channel      BIGINT UNSIGNED NOT NULL,
  rel_message      BIGINT UNSIGNED NOT NULL,
  rel_user         BIGINT UNSIGNED NOT NULL,
  rel_mentioned_by BIGINT UNSIGNED NOT NULL,

  created_at       DATETIME        NOT NULL DEFAULT NOW(),

  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE INDEX lookup_mentions ON mentions (rel_mentioned_by)
