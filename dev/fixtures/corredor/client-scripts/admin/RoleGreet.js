export default {
  label: 'Greet role (client)',
  description: 'Sits in the role editor toolbar and logs the role',

  triggers({ on }) {
    return on('manual')
      .for('system:role')
      .uiProp('app', 'admin')
      .uiProp('page', 'role/editor')
      .uiProp('slot', 'toolbar')
  },

  exec({ $role }) {
    console.info('Corredor fixture: hello role', $role.handle)
  },
}
