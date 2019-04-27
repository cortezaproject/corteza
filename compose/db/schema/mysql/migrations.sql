CREATE TABLE IF NOT EXISTS `migrations` (
 `project` varchar(16) NOT NULL COMMENT 'sam, crm, ...',
 `filename` varchar(255) NOT NULL COMMENT 'yyyymmddHHMMSS.sql',
 `statement_index` int(11) NOT NULL COMMENT 'Statement number from SQL file',
 `status` text NOT NULL COMMENT 'ok or full error message',
 PRIMARY KEY (`project`,`filename`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

