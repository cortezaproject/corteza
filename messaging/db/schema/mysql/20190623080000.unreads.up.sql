UPDATE `messaging_unread` SET rel_reply_to = 0 WHERE rel_reply_to IS NULL;
ALTER TABLE `messaging_unread` CHANGE COLUMN `rel_reply_to` `rel_reply_to` BIGINT UNSIGNED NOT NULL;
ALTER TABLE `messaging_unread` DROP PRIMARY KEY, ADD PRIMARY KEY(`rel_channel`, `rel_reply_to`, `rel_user`);

-- Add entries for all (unexisting) unreads (channels & threads)
INSERT IGNORE INTO messaging_unread
       (rel_channel, rel_reply_to, rel_user)
SELECT DISTINCT cm.rel_channel, msg.id, cm.rel_user
  FROM messaging_channel_member          AS cm
  	   INNER JOIN messaging_message AS msg ON (cm.rel_channel = msg.rel_channel AND replies > 0)
 WHERE NOT EXISTS (SELECT 1 FROM messaging_unread AS u WHERE u.rel_reply_to = msg.id AND u.rel_user = cm.rel_user)
   AND msg.rel_user > 0

UNION

SELECT DISTINCT cm.rel_channel, 0, cm.rel_user
  FROM messaging_channel_member          AS cm
 WHERE NOT EXISTS (SELECT 1 FROM messaging_unread AS u WHERE u.rel_channel = cm.rel_channel AND u.rel_user = cm.rel_user)
   AND cm.rel_user > 0
;


-- Update counters for channel messages
INSERT IGNORE INTO messaging_unread
       (rel_channel, rel_reply_to, rel_user, count, rel_last_message)
SELECT u.rel_channel, 0, u.rel_user, COUNT(m.id), u.rel_last_message
  FROM messaging_unread AS u
       INNER JOIN messaging_message AS m ON (u.rel_channel = m.rel_channel AND m.id > u.rel_last_message)
 WHERE u.rel_reply_to = 0
   AND m.reply_to = 0
 GROUP BY u.rel_channel, u.rel_user;

-- Update counters for thread messages

INSERT IGNORE INTO messaging_unread
       (rel_channel, rel_reply_to, rel_user, count, rel_last_message)
SELECT u.rel_channel, rpl.reply_to, u.rel_user, COUNT(rpl.id), u.rel_last_message
  FROM messaging_unread AS u
       INNER JOIN messaging_message AS rpl ON (u.rel_channel = rpl.rel_channel AND rpl.reply_to = u.rel_reply_to AND rpl.id > u.rel_last_message)
 WHERE rpl.replies > 0 AND u.rel_reply_to > 0
 GROUP BY u.rel_channel, rpl.reply_to, u.rel_user;
