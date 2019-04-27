CREATE TABLE crm_attachment (
  id               BIGINT UNSIGNED NOT NULL,
  rel_owner        BIGINT UNSIGNED NOT NULL,

  kind             VARCHAR(32) NOT NULL,

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

-- page attachments will be referenced via page-block meta data
-- module/record attachment will be referenced via crm_record_value
