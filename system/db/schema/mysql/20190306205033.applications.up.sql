CREATE TABLE sys_application (
  id               BIGINT UNSIGNED NOT NULL,
  rel_owner        BIGINT UNSIGNED NOT NULL REFERENCES sys_users(id),
  name             TEXT            NOT NULL COMMENT 'something we can differentiate application by',
  enabled          BOOL            NOT NULL,

  unify            JSON                NULL COMMENT 'unify specific settings',

  created_at       DATETIME        NOT NULL DEFAULT NOW(),
  updated_at       DATETIME            NULL,
  deleted_at       DATETIME            NULL, -- user soft delete

  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


REPLACE INTO `sys_application` (`id`, `name`, `enabled`, `rel_owner`, `unify`) VALUES
( 1, 'Crust Messaging', true, 0,
  '{"logo": "/applications/crust.jpg", "icon": "/applications/crust_favicon.png", "url": "/messaging/", "listed": true}'
),
( 2, 'Crust CRM', true, 0,
  '{"logo": "/applications/crust.jpg", "icon": "/applications/crust_favicon.png", "url": "/crm/", "listed": true}'
),
( 3, 'Crust Admin Area', true, 0,
  '{"logo": "/applications/crust.jpg", "icon": "/applications/crust_favicon.png", "url": "/admin/", "listed": true}'
),
( 4, 'Corteza Jitsi Bridge', true, 0,
  '{"logo": "/applications/jitsi.png", "icon": "/applications/jitsi_icon.png", "url": "/bridge/jitsi/", "listed": true}'
),
( 5, 'Google Maps', true, 0,
  '{"logo": "/applications/google_maps.png", "icon": "/applications/google_maps_icon.png", "url": "/bridge/google-maps/", "listed": true}'
);

-- Allow admin access to applications
INSERT INTO `sys_rules` (`rel_role`, `resource`, `operation`, `value`) VALUES
  (2, 'system', 'application.create', 2),
  (2, 'application:role', 'read', 2),
  (2, 'application:role', 'update', 2),
  (2, 'application:role', 'delete', 2)
;
