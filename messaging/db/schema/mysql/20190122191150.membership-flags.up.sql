ALTER TABLE channel_members ADD flag ENUM ('pinned', 'hidden', 'ignored', '') NOT NULL DEFAULT '' AFTER `type`;
