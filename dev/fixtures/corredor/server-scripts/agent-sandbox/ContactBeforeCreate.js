export default {
  label: 'Contact: default company (server)',
  description: 'Fills an empty company before a contact record is created',

  triggers({ before }) {
    return before('create')
      .for('compose:record')
      .where('module', 'agent-contact')
      .where('namespace', 'agent-sandbox')
  },

  exec({ $record }, { log }) {
    if (!$record.values.company) {
      $record.values.company = 'Set by Corredor'
      log.info('company defaulted for %s', $record.values.name)
    }

    return $record
  },
}
