DROP TABLE IF EXISTS crm_field;
DROP TABLE IF EXISTS crm_fields;
DROP TABLE IF EXISTS crm_content;
DROP TABLE IF EXISTS crm_content_links;
DROP TABLE IF EXISTS crm_content_column;
DROP TABLE IF EXISTS crm_module_content;

ALTER TABLE crm_attachment
  RENAME TO compose_attachment;

ALTER TABLE crm_chart
  RENAME TO compose_chart;

ALTER TABLE crm_module
  RENAME TO compose_module;

ALTER TABLE crm_module_form
  RENAME TO compose_module_form;

ALTER TABLE crm_page
  RENAME TO compose_page;

ALTER TABLE crm_record
  RENAME TO compose_record;

ALTER TABLE crm_record_value
  RENAME TO compose_record_value;

ALTER TABLE crm_trigger
  RENAME TO compose_trigger;
