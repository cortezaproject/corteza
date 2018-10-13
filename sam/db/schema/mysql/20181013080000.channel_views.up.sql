ALTER TABLE channel_views DROP viewed_at;
ALTER TABLE channel_views ADD rel_last_message_id BIGINT UNSIGNED;
ALTER TABLE channel_views CHANGE new_since new_messages_count INT UNSIGNED;

-- Table structure after these changes:
-- +---------------------+---------------------+------+-----+---------+-------+
-- | Field               | Type                | Null | Key | Default | Extra |
-- +---------------------+---------------------+------+-----+---------+-------+
-- | rel_channel         | bigint(20) unsigned | NO   | PRI | NULL    |       |
-- | rel_user            | bigint(20) unsigned | NO   | PRI | NULL    |       |
-- | rel_last_message_id | bigint(20) unsigned | YES  |     | NULL    |       |
-- | new_messages_count  | int(10) unsigned    | NO   |     | 0       |       |
-- +---------------------+---------------------+------+-----+---------+-------+

-- Prefill with data
INSERT INTO channel_views (rel_channel, rel_user, rel_last_message_id)
  SELECT cm.rel_channel, cm.rel_user, max(m.ID)
    FROM channel_members AS cm INNER JOIN messages AS m ON (m.rel_channel = cm.rel_channel)
  GROUP BY cm.rel_channel, cm.rel_user;

