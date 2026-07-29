// Where a revision stands in getting live, as one status the shared StatusChip
// can render. Three surfaces state this — the revision switcher's trigger, each
// entry in its dropdown, and the Publish tab's header — and they must never
// disagree about the same revision, so the derivation lives here rather than in
// any of them.
//
// It merges the two axes a revision sits on: its persisted lifecycle status and
// the session-local publish governance flag (see stores/projects.js).
import { NoID } from '@planetcrust/human-js'

// Published, running: the lifecycle statuses that mean this revision is the one
// serving users. The backend only ever sets 'active'; 'published' is accepted
// because the enum carries it.
export const LIVE_STATUSES = ['active', 'published']

// Has this revision's CHAIN ever gone live? The dashboard reports on a running
// project, so it exists only once there is something to report on — before the
// first publish the wizard is the whole product (see DashboardLayout.vue's
// redirect, RevisionSwitcher's Dashboard entry, and the list/sidebar links).
//
// Both branches read off the passed revision alone, no chain fetch needed: a
// revision with a parent proves the chain published (branching requires an
// active parent), and any status other than draft proves it directly. So only
// an original, still-unpublished draft answers false.
export function chainHasPublished(project) {
  if (!project) return false
  return project.status !== 'draft' || project.parentRevisionID !== NoID
}

export function projectStatusTag({ status, publishStatus }) {
  if (LIVE_STATUSES.includes(status)) return { status: 'active' }
  // Archived, suspended, deprecated — a revision that HAS a lifecycle verdict
  // of its own. Publish stamps the outgoing revision 'deprecated', so this is
  // most of a chain of any age, and none of it is a review state.
  if (status && status !== 'draft') return { status }
  // A draft revision is described by its review state instead — except before
  // anything is submitted, where the review axis has nothing to say and the
  // lifecycle word stands: Draft.
  return { status: publishStatus === 'draft' ? 'draft' : publishStatus }
}
