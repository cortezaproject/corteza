update channels set type = 'group' where type = 'direct';
alter table channels CHANGE type type  enum('private', 'public', 'group');
alter table channel_members CHANGE type type  enum('owner', 'member', 'invitee');
