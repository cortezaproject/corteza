// Revision work-item completeness — completed vs assigned counts (plus the
// per-status breakdown, `byStatus`: column status → its total) for ONE
// revision, off a single board call (GET /project-board/, see
// server/system/service/project_board.go). `assigned` = the sum of all four
// column totals; `completed` = the Completed column's own total. The board
// service's doc comment verified this sum is mathematically identical to the
// six-report-call count/open arithmetic the completeness bar used before —
// the single Completed status bucket IS the "not open" set. `limit: 1`
// because callers only want totals, never cards — the endpoint decouples its
// total-probe from the caller's item-page limit for exactly this caller
// shape (see loadColumn's doc comment).
//
// Shared by the two surfaces the wizard intent doc ties together: the
// Manage & Monitor rail's RevisionCompletenessBar and the publish confirm's
// unfinished-work warning (views/Wizard.vue#confirmPublish) — one function so
// the two can never drift apart on the formula.
export async function fetchRevisionCompleteness($SystemAPI, projectId, revisionId) {
  const { columns = [] } = await $SystemAPI.projectBoardBoard({
    projectID: projectId,
    revisionID: revisionId,
    limit: 1,
  })
  let assigned = 0
  let completed = 0
  const byStatus = {}
  for (const col of columns) {
    const total = Number(col.total || 0)
    assigned += total
    byStatus[col.status] = total
    if (col.status === 'Completed') completed += total
  }
  return { assigned, completed, byStatus }
}
