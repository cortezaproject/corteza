ALTER TABLE channel_views RENAME TO unreads;

ALTER TABLE unreads ADD     rel_reply_to                        BIGINT UNSIGNED NOT NULL AFTER rel_channel;
ALTER TABLE unreads CHANGE rel_channel         rel_channel      BIGINT UNSIGNED NOT NULL DEFAULT 0;
ALTER TABLE unreads CHANGE rel_user            rel_user         BIGINT UNSIGNED NOT NULL DEFAULT 0;
ALTER TABLE unreads CHANGE rel_last_message_id rel_last_message BIGINT UNSIGNED NOT NULL DEFAULT 0;
ALTER TABLE unreads CHANGE new_messages_count  count            INT    UNSIGNED NOT NULL DEFAULT 0;

