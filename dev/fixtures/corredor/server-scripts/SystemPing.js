export default {
  label: 'System ping (server)',
  description: 'Answers a manual system-level call; what a TAQ corredorExec step runs',

  triggers({ on }) {
    return on('manual').for('system')
  },

  exec(args) {
    return { pong: `pong for ${args.ping || 'nobody'}` }
  },
}
