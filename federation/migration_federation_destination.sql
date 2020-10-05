-- Pricebook module

--          id         |   rel_namespace    |  handle   |   name    | meta |          created_at           |          updated_at          | deleted_at 
-- --------------------+--------------------+-----------+-----------+------+-------------------------------+------------------------------+------------
--  196342359342989413 | 196342358990667877 | Pricebook | Pricebook | {}   | 2020-09-16 12:06:05.830318+00 | 2020-09-16 12:06:15.24715+00 | 

CREATE TABLE IF NOT EXISTS federation_module_shared
(
    id BIGINT NOT NULL,
    handle varchar(200),
    name varchar(64),
    rel_node BIGINT,
    xref_module BIGINT,
    fields TEXT,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS federation_module_exposed
(
    id BIGINT NOT NULL,
    rel_node BIGINT,
    rel_compose_module BIGINT,
    fields TEXT,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS federation_module_mapping
(
    federation_module_id BIGINT,
    compose_module_id BIGINT,
    field_mapping TEXT
);

CREATE TABLE IF NOT EXISTS federation_node
(
    id BIGINT NOT NULL,
    name varchar(256),
    status varchar(64),
    structure_status varchar(64),
    structure_synced_at varchar(64),
    data_status varchar(64),
    data_synced_at varchar(64),
    PRIMARY KEY(id)
);

INSERT INTO federation_node values ('276342359342989555', 'Node Origin (somewhere else)', 'paired', 'synced', now(), 'synced', now());

INSERT INTO federation_module_shared (id, handle, name, rel_node, xref_module, fields)
VALUES
(206342359342988000, 'Pricebook_on_origin_somewhere_else', 'Pricebook on some other origin (somewhere else)', 276342359342989555, 196342359342989000, NULL)
;
