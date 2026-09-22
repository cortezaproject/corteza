export default {
  label: 'Restricted contact script (server)',
  description:
    'Refused to super-admins and open to agent_readonly members: deny wins over allow and over the bypass roles',

  security: {
    allow: ['agent_readonly'],
    deny: ['super-admin'],
  },

  triggers({ on }) {
    return on('manual')
      .for('compose:record')
      .where('module', 'agent-contact')
      .where('namespace', 'agent-sandbox')
      .uiProp('app', 'compose')
  },

  exec({ $record }) {
    $record.values.company = 'Restricted script ran'

    return $record
  },
}
