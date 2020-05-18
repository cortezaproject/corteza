CREATE TABLE IF NOT EXISTS sys_actionlog (
  ts               DATETIME        NOT NULL DEFAULT NOW(),
  actor_ip_addr    VARCHAR(15)     NOT NULL,
  actor_id         BIGINT          UNSIGNED,
  request_origin   VARCHAR(32)     NOT NULL,
  request_id       VARCHAR(64)     NOT NULL,
  resource         VARCHAR(128)    NOT NULL,
  `action`         VARCHAR(64)     NOT NULL,
  `error`          VARCHAR(64)     NOT NULL,
  severity         SMALLINT        NOT NULL,
  description      TEXT,
  meta             JSON

) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE INDEX ts             ON sys_actionlog (ts DESC);
CREATE INDEX request_origin ON sys_actionlog (request_origin);
CREATE INDEX actor_id       ON sys_actionlog (actor_id);
CREATE INDEX resource       ON sys_actionlog (resource);
CREATE INDEX `action`       ON sys_actionlog (`action`);
