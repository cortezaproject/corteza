// Section entry gate.
//
// A section names the registry application that governs it and the user must
// hold `access` on that application; `app: null` marks the one section (home)
// everybody reaches. A section naming neither is denied — the omission has to
// fail closed, or a new section ships ungated.
//
// This is a usability boundary, not the security one: every endpoint behind it
// enforces its own permissions. What it buys is that a refused section never
// mounts, so it issues no requests and shows no half-loaded console.
//
// Both collaborators are passed in rather than imported so the rule can be
// exercised without a store, a pinia and a router around it.
export function makeSectionAccessGuard({ useApplications, sectionById, fallback = 'home' }) {
  return async to => {
    const section = sectionById(to.meta.section)
    if (section && section.app === null) return true

    const applications = useApplications()
    // The router decides before the shell has mounted, so it awaits the
    // application list itself rather than the shell's fetch of it.
    await applications.ready()

    if (section && applications.canAccessApp(section.app)) return true

    // Carry what was refused, so the landing page can say what happened
    // instead of looking like a mis-click.
    return { name: fallback, query: { denied: section?.id || to.path } }
  }
}
