export default {
  label: 'Dashboard hello (client)',
  description: 'Sits in the admin dashboard toolbar',

  triggers({ on }) {
    return on('manual')
      .for('system')
      .uiProp('app', 'admin')
      .uiProp('page', 'dashboard')
      .uiProp('slot', 'toolbar')
  },

  exec() {
    console.info('Corredor fixture: dashboard hello')
  },
}
