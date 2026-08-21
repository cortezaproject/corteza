// Naming and retention for the screenshots dev_ui_verify writes.
//
// Separate from uiverify.mjs because that file cannot be imported without a
// browser — it launches one while it loads. Everything here is plain filesystem
// work, so it is testable on its own: node --test dev/mcp/shots.test.mjs
import { readdirSync, statSync, unlinkSync } from 'node:fs'
import { join } from 'node:path'

const PREFIX = 'ui-verify-'
const SUFFIX = '.png'

// How many shots survive a sweep. Several dev-mcp servers run against one
// repo — one per agent session — so this leaves each of half a dozen recent
// sessions about ten pictures, at 40-150 KB apiece.
export const KEEP = 60

// Nothing younger than this is ever swept, whatever the count. A ui-verify run
// lasts seconds and the caller reads the path it was handed, so this is what
// stops a burst in one session from deleting a shot another session is still
// holding.
export const GRACE_MS = 10 * 60 * 1000

// One name per call, unique across concurrent callers: two live processes
// cannot share a pid, and a recycled pid cannot land in the same millisecond as
// the run that released it. A fixed name is one file per repo, and a peer
// session's page then arrives under your path.
export function shotName(now = Date.now(), pid = process.pid) {
  return `${PREFIX}${pid}-${now.toString(36)}${SUFFIX}`
}

// pruneShots keeps the newest `keep` shots in dir and deletes the rest, sparing
// anything written in the last `graceMs`. It returns the names it deleted.
//
// Only files shaped like shotName() are considered, so anything else sharing
// the directory is left alone. Every filesystem call is guarded: two sessions
// sweeping at once will each find files the other has already taken, and that
// race is the normal case here rather than a fault.
export function pruneShots(dir, { keep = KEEP, graceMs = GRACE_MS, now = Date.now() } = {}) {
  let names
  try {
    names = readdirSync(dir)
  } catch {
    return []
  }

  const shots = []
  for (const name of names) {
    if (!name.startsWith(PREFIX) || !name.endsWith(SUFFIX)) continue
    try {
      shots.push({ name, mtime: statSync(join(dir, name)).mtimeMs })
    } catch {
      continue
    }
  }

  shots.sort((a, b) => b.mtime - a.mtime)

  const deleted = []
  for (const shot of shots.slice(keep)) {
    if (now - shot.mtime < graceMs) continue
    try {
      unlinkSync(join(dir, shot.name))
      deleted.push(shot.name)
    } catch {
      continue
    }
  }

  return deleted
}
