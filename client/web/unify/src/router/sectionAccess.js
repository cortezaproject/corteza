// Section entry gate.
//
// A section names the registry application that governs it and the user must
// hold `access` on that application; `app: null` marks the one section (home)
// everybody reaches. A section naming neither is denied — the omission has to
// fail closed, or a new section ships ungated.
//
// A `perApp` section is keyed to no application of its own: it serves one
// application per route, and the one named in the route answers for itself.
//
// This is a usability boundary, not the security one: every endpoint behind it
// enforces its own permissions. What it buys is that a refused section never
// mounts, so it issues no requests and shows no half-loaded console.
//
// Both collaborators are passed in rather than imported so the rule can be
// exercised without a store, a pinia and a router around it.
// Who may enter a section. The one statement of the rule: the router asks it on
// arrival and the app menu asks it before offering a tile, so what is offered
// and what is admitted cannot drift apart.
export function sectionAllows(section, applications, application) {
  if (!section) return false
  // A per-application section is admitted by the application it is showing.
  if (section.perApp) return applicationAllows(application)
  // The ungated section (home) is reachable whatever the user holds.
  if (section.app === null) return true
  return applications.canAccessApp(section.app)
}

// A custom application speaking for itself: it exists, it is one this shell
// can render, and the user holds `access` on it.
export function applicationAllows(application) {
  return !!application && application.unify?.kind === 'custom' && !!application.canAccessApplication
}

export function makeSectionAccessGuard({ useApplications, sectionById, fallback = 'home' }) {
  return async to => {
    const section = sectionById(to.meta.section)
    if (section && section.app === null && !section.perApp) return true

    const applications = useApplications()
    // The router decides before the shell has mounted, so it awaits the
    // application list itself rather than the shell's fetch of it.
    await applications.ready()

    // Unreadable is unreachable: the read that resolves the application is the
    // same one that would put it in the menu.
    const application = section?.perApp
      ? await applications.findByID(to.params.applicationID).catch(() => undefined)
      : undefined

    if (sectionAllows(section, applications, application)) return true

    // Carry what was refused, so the landing page can say what happened
    // instead of looking like a mis-click.
    return { name: fallback, query: { denied: section?.id || to.path } }
  }
}
