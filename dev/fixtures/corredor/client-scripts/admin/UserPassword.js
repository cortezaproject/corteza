export default {
  label: 'Password hint (client)',
  description: 'Sits under the password fields of the user editor and logs the user',

  triggers({ on }) {
    return on('manual')
      .for('system:user')
      .uiProp('app', 'admin')
      .uiProp('page', 'user/editor')
      .uiProp('slot', 'passwordFooter')
  },

  exec({ $user }) {
    console.info('Corredor fixture: password slot for', $user.email)
  },
}
