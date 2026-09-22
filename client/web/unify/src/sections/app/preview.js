// A switched-off custom app is still open to whoever may change its page, as a
// preview: they need to see it before it goes to anyone else. Everybody else
// gets the disabled screen, as for any application.
export function previewsSwitchedOffApp(app) {
  return (
    !!app &&
    app.enabled === false &&
    app.unify?.kind === 'custom' &&
    !!app.canManageSourceOnApplication
  )
}
