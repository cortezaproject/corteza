export default {
  label: 'Contact: note activation (server)',
  description:
    'Appends a note to the company when a contact goes from lead to active, reading the previous values',

  triggers({ before }) {
    return before('update')
      .for('compose:record')
      .where('module', 'agent-contact')
      .where('namespace', 'agent-sandbox')
  },

  exec({ $record, $oldRecord }) {
    const was = $oldRecord && $oldRecord.values.status
    const now = $record.values.status

    if (was === 'lead' && now === 'active') {
      $record.values.company = `${$record.values.company || ''} (activated)`.trim()
    }

    return $record
  },
}
