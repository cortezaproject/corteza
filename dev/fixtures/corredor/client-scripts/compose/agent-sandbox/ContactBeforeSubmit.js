export default {
  label: 'Contact: normalise email (client)',
  description: 'Trims and lower-cases the email in the browser before the record form is submitted',

  triggers ({ before }) {
    return before('formSubmit')
      .for('ui:compose:record-page')
      .where('module', 'agent-contact')
      .where('namespace', 'agent-sandbox')
      .uiProp('app', 'compose')
  },

  exec ({ $record }) {
    const email = $record.values.email

    if (typeof email === 'string') {
      $record.values.email = email.trim().toLowerCase()
    }
  },
}
