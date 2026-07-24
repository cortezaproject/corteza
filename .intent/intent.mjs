#!/usr/bin/env node
// Intent system CLI — see .intent/SPEC.md
// Commands: check [--staged|--changed] [paths…] | sync [--all] [paths…] | status | coverage | affected <paths…>
import { createHash } from 'node:crypto'
import { execSync } from 'node:child_process'
import { readFileSync, writeFileSync, readdirSync, statSync, existsSync } from 'node:fs'
import { dirname, join, relative, resolve, basename, extname } from 'node:path'
import { fileURLToPath } from 'node:url'

const ROOT = resolve(dirname(fileURLToPath(import.meta.url)), '..')
const LOCK_PATH = join(ROOT, '.intent', 'intent.lock.json')
const config = (await import(join(ROOT, '.intent', 'config.mjs'))).default

// ---------- helpers ----------

const rel = (p) => relative(ROOT, resolve(ROOT, p)).replaceAll('\\', '/')

function globToRegex(glob) {
  let re = ''
  for (let i = 0; i < glob.length; i++) {
    const c = glob[i]
    if (c === '*') {
      if (glob[i + 1] === '*') {
        re += '.*'
        i++
        if (glob[i + 1] === '/') i++ // "**/" also matches zero dirs
      } else re += '[^/]*'
    } else if ('.+?^${}()|[]\\'.includes(c)) re += '\\' + c
    else re += c
  }
  return new RegExp('^' + re + '$')
}

const excludeRes = config.exclude.map(globToRegex)
const fileTierRes = config.fileTier.map(globToRegex)
const fileTierExclRes = (config.fileTierExclude ?? []).map(globToRegex)
const isExcluded = (p) => excludeRes.some((r) => r.test(p))
const isFileTier = (p) => fileTierRes.some((r) => r.test(p)) && !fileTierExclRes.some((r) => r.test(p))
const isSource = (p) => config.sourceExt.includes(extname(p)) && !p.endsWith('.intent.md')

function walk(dir, out) {
  for (const e of readdirSync(join(ROOT, dir), { withFileTypes: true })) {
    const p = dir + '/' + e.name
    if (isExcluded(p)) continue
    if (e.isDirectory()) walk(p, out)
    else if (isSource(p)) out.push(p)
  }
  return out
}

function enforcedFiles() {
  const out = []
  for (const root of config.enforced) {
    const abs = join(ROOT, root)
    if (!existsSync(abs)) {
      console.error(`config error: enforced root does not exist: ${root}`)
      process.exitCode = 1
      continue
    }
    if (statSync(abs).isDirectory()) walk(rel(root), out)
    else if (isSource(root) && !isExcluded(root)) out.push(rel(root))
  }
  return [...new Set(out)].sort()
}

const hashFile = (p) => createHash('sha256').update(readFileSync(join(ROOT, p))).digest('hex').slice(0, 12)

function sidecarPath(p) {
  const base = basename(p)
  return rel(join(dirname(p), base.slice(0, base.length - extname(base).length) + '.intent.md'))
}

// sidecar wins; else nearest folder doc up the tree ("." docs only govern their own dir)
function governingDoc(p) {
  const dir = dirname(p)
  const sidecar = sidecarPath(p)
  if (existsSync(join(ROOT, sidecar))) return sidecar
  let d = dir
  let hops = 0
  while (true) {
    // folder doc carries its folder's name: foo/foo.intent.md (root: INTENT.md)
    const doc = d === '.' || d === '' ? 'INTENT.md' : d + '/' + basename(d) + '.intent.md'
    if (existsSync(join(ROOT, doc)) && doc !== rel(sidecar)) {
      const fm = frontmatter(doc)
      if (hops === 0 || (fm.covers ?? 'recursive') === 'recursive') return rel(doc)
    }
    if (d === '.' || d === '') return null
    d = dirname(d)
    hops++
  }
}

// minimal YAML-subset frontmatter parser: scalars, "- item" lists, inline comments
function frontmatter(docPath) {
  const text = readFileSync(join(ROOT, docPath), 'utf8')
  if (!text.startsWith('---\n')) return { __error: 'missing frontmatter' }
  const end = text.indexOf('\n---', 4)
  if (end === -1) return { __error: 'unterminated frontmatter' }
  const out = { __body: text.slice(end + 4) }
  let key = null
  for (const raw of text.slice(4, end).split('\n')) {
    const line = raw.replace(/(^|\s)#.*$/, '').trimEnd()
    if (!line.trim()) continue
    if (/^\s+- /.test(line)) {
      if (!key) return { __error: `list item without key: ${line.trim()}` }
      out[key].push(line.trim().slice(2).trim())
    } else {
      const m = line.match(/^([\w-]+):\s*(.*)$/)
      if (!m) return { __error: `unparseable line: ${line.trim()}` }
      key = m[1]
      out[key] = m[2] === '' || m[2] === '[]' ? [] : m[2].replace(/^["']|["']$/g, '')
    }
  }
  return out
}

function loadLock() {
  return existsSync(LOCK_PATH) ? JSON.parse(readFileSync(LOCK_PATH, 'utf8')) : { version: 1, files: {} }
}
function saveLock(lock) {
  lock.files = Object.fromEntries(Object.entries(lock.files).sort(([a], [b]) => a.localeCompare(b)))
  writeFileSync(LOCK_PATH, JSON.stringify(lock, null, 2) + '\n')
}

function gitList(args) {
  try {
    return execSync(`git ${args}`, { cwd: ROOT, encoding: 'utf8' }).split('\n').filter(Boolean).map(rel)
  } catch {
    return []
  }
}

// ---------- doc validation ----------

function findDocs() {
  const out = []
  const walkDocs = (dir) => {
    for (const e of readdirSync(join(ROOT, dir === '.' ? '' : dir), { withFileTypes: true })) {
      const p = (dir === '.' ? '' : dir + '/') + e.name
      if (isExcluded(p) || p.startsWith('.intent/') || p.startsWith('.git/')) continue
      if (e.isDirectory()) walkDocs(p)
      else if (e.name.endsWith('.intent.md')) out.push(p)
    }
  }
  for (const root of config.enforced) {
    const abs = join(ROOT, root)
    if (existsSync(abs) && statSync(abs).isDirectory()) walkDocs(rel(root))
  }
  if (existsSync(join(ROOT, 'INTENT.md'))) out.push('INTENT.md')
  return out
}

function validateDoc(p) {
  const errs = []
  const fm = frontmatter(p)
  if (fm.__error) return [`${p}: frontmatter — ${fm.__error}`]
  if (!['folder', 'file'].includes(fm.kind)) errs.push(`${p}: kind must be folder|file`)
  if (fm.kind === 'file') {
    const covers = join(dirname(p), fm.covers ?? '')
    if (!fm.covers || !existsSync(join(ROOT, covers))) errs.push(`${p}: covers must name an existing sibling file`)
  }
  const prose = fm.__body.split('\n').filter((l) => l.trim()).length
  if (prose > config.docMaxLines) errs.push(`${p}: ${prose} prose lines (cap ${config.docMaxLines})`)
  return errs
}

// ---------- commands ----------

function computeDrift(files, lock) {
  const drift = []
  for (const f of files) {
    const entry = lock.files[f]
    const h = hashFile(f)
    if (!entry) drift.push({ file: f, why: 'never synced', doc: governingDoc(f) })
    else if (entry.h !== h) drift.push({ file: f, why: 'changed since sync', doc: governingDoc(f) })
  }
  const orphans = Object.keys(lock.files).filter((f) => !existsSync(join(ROOT, f)))
  return { drift, orphans }
}

function cmdCheck(argv) {
  const files = enforcedFiles()
  let scope = files
  if (argv.includes('--staged')) {
    const staged = new Set(gitList('diff --cached --name-only'))
    scope = files.filter((f) => staged.has(f))
  } else if (argv.includes('--changed')) {
    const changed = new Set([...gitList('diff --name-only'), ...gitList('diff --cached --name-only'), ...gitList('ls-files --others --exclude-standard')])
    scope = files.filter((f) => changed.has(f))
  }
  const lock = loadLock()
  const { drift, orphans } = computeDrift(scope, lock)
  const docErrs = findDocs().flatMap(validateDoc)
  const uncovered = scope.filter((f) => !governingDoc(f))

  let failed = false
  if (drift.length) {
    failed = true
    console.error('DRIFT — covered files changed without their intent doc being re-synced:')
    for (const d of drift) console.error(`  ${d.file}  (${d.why})  → update & sync: ${d.doc ?? 'NO GOVERNING DOC'}`)
  }
  if (!argv.includes('--staged') && orphans.length) {
    failed = true
    console.error('ORPHANED lock entries (file moved/deleted — re-sync to acknowledge):')
    for (const o of orphans) console.error(`  ${o}`)
  }
  if (docErrs.length) {
    failed = true
    console.error('DOC ERRORS:')
    for (const e of docErrs) console.error(`  ${e}`)
  }
  if (uncovered.length) {
    failed = true
    console.error('NO GOVERNING DOC (coverage gap):')
    for (const f of uncovered) console.error(`  ${f}`)
  }
  if (failed) {
    console.error('\nintent check FAILED. Reconcile the intent docs, then: node .intent/intent.mjs sync')
    process.exit(1)
  }
  console.log(`intent check OK (${scope.length} files, ${findDocs().length} docs)`)
}

function cmdSync(argv) {
  const files = enforcedFiles()
  const paths = argv.filter((a) => !a.startsWith('--')).map(rel)
  const scope = paths.length ? files.filter((f) => paths.some((p) => f === p || f.startsWith(p + '/'))) : files
  const lock = loadLock()
  const now = new Date().toISOString().slice(0, 10)
  let synced = 0
  const byDoc = new Map()
  for (const f of scope) {
    const doc = governingDoc(f)
    if (!doc) {
      console.error(`SKIP (no governing doc): ${f}`)
      continue
    }
    const h = hashFile(f)
    if (lock.files[f]?.h === h && lock.files[f]?.doc === doc) continue
    lock.files[f] = { h, doc, t: now }
    synced++
    byDoc.set(doc, (byDoc.get(doc) ?? 0) + 1)
  }
  for (const o of Object.keys(lock.files).filter((f) => !existsSync(join(ROOT, f)))) delete lock.files[o]
  saveLock(lock)
  for (const [doc, n] of byDoc) console.log(`  synced ${n} file(s) governed by ${doc}`)
  console.log(`intent sync: ${synced} file(s) updated in lock`)
}

function cmdStatus() {
  const files = enforcedFiles()
  const lock = loadLock()
  const { drift, orphans } = computeDrift(files, lock)
  console.log(`enforced roots: ${config.enforced.length}, covered files: ${files.length}, docs: ${findDocs().length}`)
  console.log(`drift: ${drift.length}, orphaned lock entries: ${orphans.length}`)
  const byRoot = new Map()
  for (const root of config.enforced) {
    const rooted = files.filter((f) => f === rel(root) || f.startsWith(rel(root) + '/'))
    const oldest = rooted.map((f) => lock.files[f]?.t).filter(Boolean).sort()[0]
    byRoot.set(root, oldest ?? 'never')
  }
  console.log('drift-age (oldest sync per enforced root — audit oldest first):')
  for (const [root, t] of [...byRoot].sort((a, b) => String(a[1]).localeCompare(String(b[1])))) console.log(`  ${t}  ${root}`)
  const pending = config.covered.filter((c) => !config.enforced.some((e) => c === e || c.startsWith(e + '/') || e.startsWith(c + '/')))
  if (pending.length) console.log(`backfill pending (covered, not yet enforced): ${pending.join(', ')}`)
}

function cmdCoverage() {
  const files = enforcedFiles()
  const uncovered = files.filter((f) => !governingDoc(f))
  // file-tier paths must be governed by their own sidecar, not just a folder doc
  const missingSidecar = files.filter((f) => isFileTier(f) && governingDoc(f) !== sidecarPath(f))
  if (missingSidecar.length) {
    console.error('file-tier files MISSING their <name>.intent.md sidecar:')
    for (const f of missingSidecar) console.error(`  ${f}`)
  }
  const owned = new Set(files.map((f) => governingDoc(f)).filter(Boolean))
  const orphanDocs = findDocs().filter((d) => !owned.has(d))
  if (uncovered.length) {
    console.error('files with NO governing doc:')
    for (const f of uncovered) console.error(`  ${f}`)
  }
  // Folder docs legitimately govern nothing directly when every child carries
  // its own doc — count them. File-kind orphans are real problems and print.
  // (Moved/deleted sources are the lockfile's job, not this note.)
  const unexpected = orphanDocs.filter((d) => frontmatter(d).kind !== 'folder')
  for (const d of unexpected) console.log(`  note: orphaned file doc: ${d}`)
  const idx = orphanDocs.length - unexpected.length
  if (idx) console.log(`  (${idx} folder docs govern no files directly — expected)`)
  if (uncovered.length || missingSidecar.length) process.exit(1)
  console.log(`intent coverage OK (${files.length} files all governed)`)
}

// Tests to run for a set of changed files:
// 1. the governing doc's `tests` for each changed file, plus
// 2. `tests` of any doc that declares a changed file (or its folder) in
//    `depends-on` — consumers whose contract may be affected.
function cmdAffected(argv) {
  const changed = argv.filter((a) => !a.startsWith('--')).map(rel)
  const tests = new Set()
  const addTests = (fm) => {
    for (const t of Array.isArray(fm.tests) ? fm.tests : []) tests.add(t)
  }
  for (const p of changed) {
    const doc = existsSync(join(ROOT, p)) && governingDoc(p)
    if (doc) addTests(frontmatter(doc))
  }
  for (const doc of findDocs()) {
    const fm = frontmatter(doc)
    const deps = Array.isArray(fm['depends-on']) ? fm['depends-on'] : []
    if (deps.some((d) => changed.some((c) => c === d || c.startsWith(d + '/')))) addTests(fm)
  }
  for (const t of [...tests].sort()) console.log(t)
  if (!tests.size) console.error('(no mapped tests for the given files)')
}

const [cmd, ...argv] = process.argv.slice(2)
const commands = { check: cmdCheck, sync: cmdSync, status: cmdStatus, coverage: cmdCoverage, affected: cmdAffected }
if (!commands[cmd]) {
  console.error('usage: intent.mjs <check|sync|status|coverage|affected> [args]')
  process.exit(2)
}
commands[cmd](argv)
