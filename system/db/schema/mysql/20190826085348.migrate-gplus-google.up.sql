/* migrates existing credentials */
UPDATE sys_credentials SET kind = 'google' WHERE kind = 'gplus';

/* migrates existing settings. */
UPDATE sys_settings SET name = REPLACE(name, '.gplus.', '.google.') WHERE name LIKE 'auth.external.providers.gplus.%';
