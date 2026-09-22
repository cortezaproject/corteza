export default {
  label: 'Contact: normalise email on admin pages (client)',
  description: 'Trims and lower-cases the email before the admin record form is submitted',

  triggers({ before }) {
    return before('formSubmit')
      .for('ui:compose:admin-record-page')
      .where('module', 'agent-contact')
      .where('namespace', 'agent-sandbox')
      .uiProp('app', 'compose')
  },

  exec({ $record }) {
    const email = $record.values.email

    if (typeof email === 'string') {
      $record.values.email = email.trim().toLowerCase()
    }
  },
}
