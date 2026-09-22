export default {
  label: 'Task: heartbeat (server, interval)',
  description:
    'Keeps one heartbeat task stamped with the last minute it ran, as the run-as user; proves deferred triggers and run-as',

  security: {
    runAs: 'agent@local.dev',
  },

  triggers({ every }) {
    return every('* * * * *')
  },

  async exec(args, { Compose, log }) {
    const ns = await Compose.resolveNamespace('agent-sandbox')
    const module = await Compose.findModuleByHandle('agent-task', ns)
    const { set = [] } = await Compose.findRecords("title = 'Corredor heartbeat'", module)

    const record = set[0] || (await Compose.makeRecord({ title: 'Corredor heartbeat' }, module))
    record.values.due = new Date().toISOString()

    const saved = await Compose.saveRecord(record)
    log.info('heartbeat task %s stamped', saved.recordID)
  },
}
