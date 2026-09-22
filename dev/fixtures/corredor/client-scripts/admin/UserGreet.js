export default {
  label: 'Greet user (client)',
  description: 'Logs the user being edited to the browser console',

  triggers ({ on }) {
    return on('manual')
      .for('system:user')
      .uiProp('app', 'admin')
      .uiProp('page', 'user/editor')
      .uiProp('slot', 'infoFooter')
  },

  exec ({ $user }) {
    console.info('Corredor fixture: hello', $user.email)
  },
}
