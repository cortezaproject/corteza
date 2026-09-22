export default {
  label: 'Activate contact (server)',
  description: 'Marks the contact active; runs inside Corredor when its page button is clicked',

  triggers ({ on }) {
    return on('manual')
      .for('compose:record')
      .where('module', 'agent-contact')
      .where('namespace', 'agent-sandbox')
      .uiProp('app', 'compose')
  },

  exec ({ $record }) {
    $record.values.status = 'active'

    return $record
  },
}
