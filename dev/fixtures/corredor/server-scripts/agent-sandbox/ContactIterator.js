export default {
  label: 'Contact: sweep dormant (server, iterator)',
  description:
    'Every minute, walks the dormant contacts and stamps their company; proves iterator scripts',

  security: {
    runAs: 'agent@local.dev',
  },

  iterator(each) {
    return each({
      resourceType: 'compose:record',
      action: 'update',
      filter: {
        namespace: 'agent-sandbox',
        module: 'agent-contact',
        query: "status = 'dormant'",
        limit: 10,
      },
    }).every('* * * * *')
  },

  exec({ $record }) {
    $record.values.company = 'Swept by Corredor'

    return $record
  },
}
