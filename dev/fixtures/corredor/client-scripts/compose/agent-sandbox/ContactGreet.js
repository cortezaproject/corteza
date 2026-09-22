export default {
  label: 'Greet contact (client)',
  description: 'Shows a toast in the browser and opens the record viewer',

  triggers ({ on }) {
    return on('manual')
      .for('compose:record')
      .where('module', 'agent-contact')
      .where('namespace', 'agent-sandbox')
      .uiProp('app', 'compose')
  },

  exec ({ $record }, { ComposeUI }) {
    ComposeUI.success(`Hello, ${$record.values.name || 'stranger'} (client script)`)
    ComposeUI.gotoRecordViewer()
  },
}
