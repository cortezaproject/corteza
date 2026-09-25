// Home entry redirect.
//
// Opening the shell at a bare `/` goes to the home application instead of the
// launcher: the user's own pick (`meta.homeApplicationID`) or else the one the
// instance marks `unify.home`. Only the first navigation of a page load counts,
// so the home button, a denied bounce and any other in-app trip to `/` still
// show the launcher.
//
// A home application the user cannot open leaves them on the launcher, told
// so. Collaborators are passed in so the rule runs without a router or store.
import { START_LOCATION } from 'vue-router'

const NO_ID = '0'

// The application `/` should open, or why it cannot: { app } | { unavailable }
// | {} when nothing is set.
export function resolveHome({ ownID, apps, openable }) {
  const own = ownID && ownID !== NO_ID
  const app = own ? apps.find(a => a.applicationID === ownID) : apps.find(a => a.unify?.home)

  if (!app) return own ? { unavailable: '' } : {}
  if (!openable(app)) return { unavailable: app.unify?.name || app.name || '' }
  return { app }
}

export function makeHomeEntryGuard({ useApplications, ownHomeID, openable, hrefOf, leave }) {
  return async (to, from) => {
    if (from !== START_LOCATION || to.name !== 'home' || Object.keys(to.query).length) return true

    const applications = useApplications()
    await applications.ready()

    const { app, unavailable } = resolveHome({
      ownID: ownHomeID(),
      apps: applications.apps,
      openable,
    })

    if (unavailable !== undefined) return { name: 'home', query: { homeUnavailable: unavailable } }
    if (!app) return true

    const { href, local } = hrefOf(app)
    if (local) return href

    leave(href)
    return false
  }
}
