export default {
  label: 'Greet user group (client)',
  description: 'Sits under the basic information of the user group editor',

  triggers({ on }) {
    return on('manual')
      .for('system:user-group')
      .uiProp('app', 'admin')
      .uiProp('page', 'user-group/editor')
      .uiProp('slot', 'infoFooter')
  },

  exec({ $userGroup }) {
    console.info('Corredor fixture: hello group', $userGroup && $userGroup.name)
  },
}
