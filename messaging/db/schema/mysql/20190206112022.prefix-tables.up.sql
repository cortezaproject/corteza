-- misc tables

ALTER TABLE attachments            RENAME TO messaging_attachment;
ALTER TABLE mentions               RENAME TO messaging_mention;
ALTER TABLE unreads                RENAME TO messaging_unread;

-- channel tables

ALTER TABLE channels               RENAME TO messaging_channel;
ALTER TABLE channel_members        RENAME TO messaging_channel_member;

-- message tables

ALTER TABLE messages               RENAME TO messaging_message;
ALTER TABLE message_attachment     RENAME TO messaging_message_attachment;
ALTER TABLE message_flags          RENAME TO messaging_message_flag;
