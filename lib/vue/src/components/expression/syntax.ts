// Tokenising and checking for the expression inputs.
//
// These are string functions with no editor dependency, so the rules can be
// tested directly; CInputExpression adapts their output to CodeMirror.

import { membersAt, resolvePath, type ScopeEntry } from './catalog'

// Keywords the QL lexer reserves — server/pkg/ql/token_consumers.go.
export const QL_KEYWORDS = [
  'AND',
  'OR',
  'XOR',
  'NOT',
  'IS',
  'IN',
  'LIKE',
  'BETWEEN',
  'NULL',
  'TRUE',
  'FALSE',
]

// Functions the server-evaluated expression language registers — mirrors
// server/pkg/expr/func_names.go, whose TestBuiltInFunctionNames proves each one
// is really reserved. Calling a name that is not here is a hard evaluation
// error, not a silent null.
export const EXPR_FUNCTIONS = [
  'abs',
  'average',
  'base64encode',
  'camelize',
  'ceil',
  'coalesce',
  'count',
  'earliest',
  'filter',
  'find',
  'float',
  'floor',
  'format',
  'has',
  'hasAll',
  'hasPrefix',
  'hasSubstring',
  'hasSuffix',
  'int',
  'isEmail',
  'isEmpty',
  'isLeapYear',
  'isNil',
  'isUrl',
  'isWeekDay',
  'join',
  'latest',
  'length',
  'log',
  'longest',
  'match',
  'max',
  'merge',
  'min',
  'modDate',
  'modMonth',
  'modTime',
  'modWeek',
  'modYear',
  'now',
  'omit',
  'parseDuration',
  'parseISOTime',
  'pop',
  'pow',
  'push',
  'random',
  'round',
  'set',
  'shift',
  'shortest',
  'snakify',
  'sort',
  'splice',
  'split',
  'sqrt',
  'strftime',
  'sub',
  'sum',
  'title',
  'toJSON',
  'toLower',
  'toUpper',
  'trim',
  'trimLeft',
  'trimRight',
  'untitle',
]

// Literals the language accepts as bare words; never variables.
const EXPR_LITERALS = ['true', 'false', 'null', 'nil']

export type Dialect = 'ql' | 'interpolation' | 'expr'

export interface Hole {
  // Offset of the opening `${`.
  from: number
  // Offset just past the closing `}`, or the end of the text when unterminated.
  to: number
  // Text between the braces.
  inner: string
  terminated: boolean
}

// Finds the `${...}` holes in an author-written template.
//
// Brace depth is tracked so a hole containing an object literal closes at the
// right place, and quoted runs are skipped so a brace inside a string does not
// move that depth.
export function scanHoles(text: string): Hole[] {
  const holes: Hole[] = []

  for (let i = 0; i < text.length - 1; i++) {
    if (text[i] !== '$' || text[i + 1] !== '{') continue

    let depth = 1
    let quote = ''
    let j = i + 2

    for (; j < text.length; j++) {
      const ch = text[j]

      if (quote) {
        if (ch === '\\') j++
        else if (ch === quote) quote = ''
        continue
      }

      if (ch === "'" || ch === '"' || ch === '`') quote = ch
      else if (ch === '{') depth++
      else if (ch === '}') {
        depth--
        if (depth === 0) break
      }
    }

    const terminated = depth === 0 && j < text.length
    const to = terminated ? j + 1 : text.length

    holes.push({ from: i, to, inner: text.slice(i + 2, terminated ? j : text.length), terminated })

    i = to - 1
  }

  return holes
}

// A hole whose whole content is a dotted path, and so can be checked against
// the scope. Anything richer — `${a || 'b'}`, a call, an index — is a valid
// template that this module deliberately does not try to understand, because
// reporting it as wrong would be worse than saying nothing.
const DOTTED_PATH = /^\s*([A-Za-z_$][\w$]*(?:\s*\.\s*[A-Za-z_$][\w$]*)*)\s*$/

export interface Diagnostic {
  from: number
  to: number
  severity: 'error' | 'warning'
  message: string
}

// Identifier paths in a server-evaluated expression, with what precedes and
// follows each, so a function call and a member access can be told apart.
export interface IdentRef {
  from: number
  to: number
  path: string[]
  isCall: boolean
}

const IDENT_PATH = /[A-Za-z_$][\w$]*(?:\s*\.\s*[A-Za-z_$][\w$]*)*/g

// Scans identifier paths outside string literals. Deliberately not a parser:
// it finds the names an expression mentions, which is all the checks below and
// the completions need.
export function scanIdents(text: string): IdentRef[] {
  const strings = tokenRanges(text, 'expr').filter(r => r.kind === 'string')
  const inString = (i: number) => strings.some(r => i >= r.from && i < r.to)
  const out: IdentRef[] = []

  IDENT_PATH.lastIndex = 0
  let m: RegExpExecArray | null

  while ((m = IDENT_PATH.exec(text))) {
    if (inString(m.index)) continue

    const to = m.index + m[0].length
    // A member access reached through a dot is part of its own path already;
    // only a name with nothing but whitespace and `(` after it is a call.
    const isCall = /^\s*\(/.test(text.slice(to))

    out.push({
      from: m.index,
      to,
      path: m[0].split('.').map(s => s.trim()),
      isCall,
    })
  }

  return out
}

// Checks a server-evaluated expression.
//
// Severity follows what the evaluator actually does, measured against
// /system/expressions/evaluate: an unknown root or an unknown function is a
// hard error that fails the whole call — and `usePageVisibility` treats a
// failed call as "every expression false", so the block silently disappears.
// An unknown *member* of a known root merely resolves to null, so it is a
// warning: wrong, almost certainly, but not destructive.
function lintExprDialect(text: string, scope: ScopeEntry[]): Diagnostic[] {
  const out: Diagnostic[] = []

  for (const ref of scanIdents(text)) {
    if (ref.isCall) {
      const name = ref.path[ref.path.length - 1]
      if (ref.path.length === 1 && !EXPR_FUNCTIONS.includes(name)) {
        out.push({
          from: ref.from,
          to: ref.to,
          severity: 'error',
          message: `'${name}' is not a known function — the whole expression fails to evaluate`,
        })
      }
      continue
    }

    if (EXPR_LITERALS.includes(ref.path[0])) continue

    const res = resolvePath(scope, ref.path)
    if (res.status !== 'unknown') continue

    if (!res.parent) {
      const roots = scope.filter(e => e.suggest !== false).map(e => e.name)
      out.push({
        from: ref.from,
        to: ref.to,
        severity: 'error',
        message:
          `'${res.segment}' is not an available variable — the whole expression fails to ` +
          `evaluate, which hides this${roots.length ? `. Available: ${roots.join(', ')}` : ''}`,
      })
      continue
    }

    const known = (res.parent.fields || []).filter(e => e.suggest !== false).map(e => e.name)
    out.push({
      from: ref.from,
      to: ref.to,
      severity: 'warning',
      message: `'${res.segment}' is not a member of ${res.parent.name}, so this reads as empty${
        known.length ? ` — try ${known.slice(0, 6).join(', ')}` : ''
      }`,
    })
  }

  return out
}

export function lintExpression(
  text: string,
  scope: ScopeEntry[],
  dialect: Dialect = 'ql',
): Diagnostic[] {
  if (dialect === 'expr') return lintExprDialect(text, scope)

  const out: Diagnostic[] = []

  for (const hole of scanHoles(text)) {
    if (!hole.terminated) {
      out.push({
        from: hole.from,
        to: hole.to,
        severity: 'error',
        message: 'Unterminated ${ — add a closing }',
      })
      continue
    }

    if (!hole.inner.trim()) {
      out.push({ from: hole.from, to: hole.to, severity: 'error', message: 'Empty ${}' })
      continue
    }

    const match = DOTTED_PATH.exec(hole.inner)
    if (!match) continue

    const path = match[1].split('.').map(s => s.trim())
    const res = resolvePath(scope, path)
    if (res.status !== 'unknown') continue

    const known = (res.parent ? res.parent.fields || [] : scope)
      .filter(e => e.suggest !== false)
      .map(e => e.name)

    out.push({
      from: hole.from,
      to: hole.to,
      severity: 'error',
      message: res.parent
        ? `'${res.segment}' is not a member of ${res.parent.name}${
            known.length ? ` — try ${known.slice(0, 6).join(', ')}` : ''
          }`
        : `'${res.segment}' is not an available variable${
            known.length ? ` — try ${known.join(', ')}` : ''
          }`,
    })
  }

  return out
}

export interface TokenRange {
  from: number
  to: number
  kind: 'hole' | 'keyword' | 'string' | 'number' | 'function'
}

const QL_STRING = /'(?:\\.|[^'\\])*'/g
const EXPR_STRING = /'(?:\\.|[^'\\])*'|"(?:\\.|[^"\\])*"/g
const EXPR_LITERAL = new RegExp(`\\b(?:${EXPR_LITERALS.join('|')})\\b`, 'g')
const EXPR_CALL = new RegExp(`\\b(?:${EXPR_FUNCTIONS.join('|')})\\b(?=\\s*\\()`, 'g')
const QL_NUMBER = /\b\d+(?:\.\d+)?\b/g
const QL_KEYWORD = new RegExp(`\\b(?:${QL_KEYWORDS.join('|')})\\b`, 'gi')

// Ranges to colour, in document order and never overlapping.
//
// Holes are found first and QL tokens only in the gaps between them, so a
// keyword inside `${...}` stays part of the hole rather than being painted
// twice — CodeMirror rejects overlapping ranges from a single builder.
export function tokenRanges(text: string, dialect: Dialect): TokenRange[] {
  // A server-evaluated expression has no `${}` holes — it is the expression.
  if (dialect === 'expr') {
    const out: TokenRange[] = []

    for (const [re, kind] of [
      [EXPR_STRING, 'string'],
      [EXPR_CALL, 'function'],
      [EXPR_LITERAL, 'keyword'],
      [QL_NUMBER, 'number'],
    ] as Array<[RegExp, TokenRange['kind']]>) {
      re.lastIndex = 0
      let m: RegExpExecArray | null
      while ((m = re.exec(text))) {
        out.push({ from: m.index, to: m.index + m[0].length, kind })
      }
    }

    return disjoint(out)
  }

  const holes = scanHoles(text)
  const out: TokenRange[] = holes.map(h => ({ from: h.from, to: h.to, kind: 'hole' as const }))

  if (dialect === 'ql') {
    let cursor = 0

    for (const gap of [
      ...holes.map(h => ({ from: h.from, to: h.to })),
      { from: text.length, to: text.length },
    ]) {
      const slice = text.slice(cursor, gap.from)
      const offset = cursor

      for (const [re, kind] of [
        [QL_STRING, 'string'],
        [QL_KEYWORD, 'keyword'],
        [QL_NUMBER, 'number'],
      ] as Array<[RegExp, TokenRange['kind']]>) {
        re.lastIndex = 0
        let m: RegExpExecArray | null
        while ((m = re.exec(slice))) {
          out.push({ from: offset + m.index, to: offset + m.index + m[0].length, kind })
        }
      }

      cursor = gap.to
    }
  }

  return disjoint(out)
}

// A number or keyword inside a string literal would overlap it; strings are
// pushed first, so dropping later ranges that start inside an earlier one keeps
// the set disjoint — CodeMirror rejects overlaps from a single builder.
function disjoint(ranges: TokenRange[]): TokenRange[] {
  const sorted = [...ranges].sort((a, b) => a.from - b.from || b.to - a.to)
  const out: TokenRange[] = []
  let end = -1

  for (const r of sorted) {
    if (r.from < end) continue
    out.push(r)
    end = r.to
  }

  return out
}

// Record columns a prefilter may name directly, alongside the module's own
// fields. These are the DAL model's system attributes
// (server/compose/service/module.go) — `recordID` is an accepted alias for the
// primary `ID`. Storage-infrastructure attributes the server also resolves
// (tenantID, projectID, createdByAgent) are left out: valid to type, but not
// something a page author filters on.
export const QL_SYSTEM_FIELDS = [
  'recordID',
  'ownedBy',
  'createdAt',
  'createdBy',
  'updatedAt',
  'updatedBy',
  'deletedAt',
  'deletedBy',
]

// Ranking between the groups an author picks from. What they are filtering
// lands above the record's own bookkeeping, and language keywords last.
const BOOST = { field: 50, system: 25, keyword: 0 }

export interface CompletionOption {
  label: string
  detail?: string
  // Text written in place of [from, to). Defaults to `label`; differs where the
  // completion carries punctuation, e.g. wrapping a variable in `${…}`.
  insert?: string
  // Caret offset within `insert`. Defaults to its end.
  cursor?: number
  // Reopen the completion list after applying — used where the accepted option
  // leaves the author mid-path, e.g. `${record.`.
  retrigger?: boolean
  // Higher sorts first.
  boost?: number
}

export interface CompletionResult {
  from: number
  to: number
  options: CompletionOption[]
}

function detailOf(entry: ScopeEntry): string {
  return entry.label && entry.label !== entry.name ? `${entry.type} · ${entry.label}` : entry.type
}

// The hole the cursor sits inside, if any.
function holeAt(text: string, pos: number): Hole | null {
  return (
    scanHoles(text).find(h => pos > h.from + 1 && pos <= (h.terminated ? h.to - 1 : h.to)) || null
  )
}

// Whether pos sits inside a QL string literal, where a `$` is currency rather
// than the start of a variable.
//
// Counts quotes rather than reusing tokenRanges: the string being typed has no
// closing quote yet, which is precisely when this is asked. Quotes inside a
// `${…}` hole belong to the hole's own expression and are skipped.
function inStringLiteral(text: string, pos: number, dialect: Dialect): boolean {
  if (dialect !== 'ql') return false

  const holes = scanHoles(text.slice(0, pos))
  let open = false

  for (let i = 0; i < pos; i++) {
    const hole = holes.find(h => i >= h.from && i < h.to)
    if (hole) {
      i = hole.to - 1
      continue
    }
    if (text[i] === '\\') i++
    else if (text[i] === "'") open = !open
  }

  return open
}

const startsWith = (name: string, word: string) =>
  !word || name.toLowerCase().startsWith(word.toLowerCase())

// Highest boost first, order within a group preserved — a module's fields stay
// in the order their author arranged them. The editor consumes this order as
// given rather than re-sorting.
function ranked(options: CompletionOption[]): CompletionOption[] {
  return options
    .map((o, i) => ({ o, i }))
    .sort((a, b) => (b.o.boost ?? 0) - (a.o.boost ?? 0) || a.i - b.i)
    .map(({ o }) => o)
}

// What to offer at `pos`.
//
// Three places can complete: inside a `${…}` hole, immediately after a bare `$`
// (which offers the same scope but writes the braces), and — in QL only — on a
// bare identifier, where the module's fields and the record's system columns
// are named. QL identifiers are never reported as *wrong*, since the language's
// own functions and literals are indistinguishable from a mistyped field here.
//
// `explicit` is set when the author asked for the list (Ctrl-Space) rather than
// typing into it; then an empty prefix offers everything instead of nothing.
export function completionAt(
  text: string,
  pos: number,
  scope: ScopeEntry[],
  dialect: Dialect,
  queryFields: Array<{ name: string; label?: string; kind?: string }> = [],
  explicit = false,
): CompletionResult | null {
  // A server-evaluated expression is a bare expression: the scope and the
  // language's functions are offered directly, with no `${}` to open first.
  if (dialect === 'expr') {
    const typed = text.slice(0, pos)
    const pathMatch = /([A-Za-z_$][\w$]*(?:\.[A-Za-z_$][\w$]*)*\.)?([A-Za-z_$][\w$]*)?$/.exec(typed)
    if (!pathMatch) return null

    const prefixPath = (pathMatch[1] || '').split('.').filter(Boolean)
    const word = pathMatch[2] || ''
    if (!word && !prefixPath.length && !explicit) return null

    // Inside a path, only that object's members make sense.
    if (prefixPath.length) {
      const members = membersAt(scope, prefixPath).filter(e => startsWith(e.name, word))
      return members.length
        ? {
            from: pos - word.length,
            to: pos,
            options: ranked(
              members.map(e => ({
                label: e.name,
                detail: detailOf(e),
                insert: e.fields ? `${e.name}.` : e.name,
                retrigger: !!e.fields,
                boost: e.fields ? BOOST.system : BOOST.field,
              })),
            ),
          }
        : null
    }

    const options: CompletionOption[] = [
      ...membersAt(scope, [])
        .filter(e => startsWith(e.name, word))
        .map(e => ({
          label: e.name,
          detail: detailOf(e),
          insert: e.fields ? `${e.name}.` : e.name,
          retrigger: !!e.fields,
          boost: BOOST.field,
        })),
      ...EXPR_FUNCTIONS.filter(n => startsWith(n, word)).map(n => ({
        label: n,
        detail: 'function',
        // Land between the parentheses, ready for the argument.
        insert: `${n}()`,
        cursor: n.length + 1,
        boost: BOOST.system,
      })),
    ]

    return options.length ? { from: pos - word.length, to: pos, options: ranked(options) } : null
  }

  const hole = holeAt(text, pos)

  if (hole) {
    const typed = text.slice(hole.from + 2, pos)
    const pathMatch = /([A-Za-z_$][\w$]*(?:\.[A-Za-z_$][\w$]*)*\.)?([A-Za-z_$][\w$]*)?$/.exec(typed)
    if (!pathMatch) return null

    const prefixPath = (pathMatch[1] || '').split('.').filter(Boolean)
    const word = pathMatch[2] || ''
    const members = membersAt(scope, prefixPath)
    if (!members.length) return null

    const options = members
      .filter(e => startsWith(e.name, word))
      .map(e => ({
        label: e.name,
        detail: detailOf(e),
        // Stepping into an object leaves the author mid-path, so reopen the
        // list on the members they just asked for.
        insert: e.fields ? `${e.name}.` : e.name,
        retrigger: !!e.fields,
        boost: e.fields ? BOOST.system : BOOST.field,
      }))

    return options.length ? { from: pos - word.length, to: pos, options: ranked(options) } : null
  }

  // A bare `$` — offer the scope and write the braces around the choice.
  const dollar = /\$([A-Za-z_$][\w$]*)?$/.exec(text.slice(0, pos))
  if (dollar && !inStringLiteral(text, pos, dialect)) {
    const word = dollar[1] || ''
    const options = membersAt(scope, [])
      .filter(e => startsWith(e.name, word))
      .map(e => ({
        label: e.name,
        detail: detailOf(e),
        insert: e.fields ? `\${${e.name}.}` : `\${${e.name}}`,
        // Land inside the braces when there is more path to type.
        cursor: e.fields ? e.name.length + 3 : undefined,
        retrigger: !!e.fields,
        boost: e.fields ? BOOST.system : BOOST.field,
      }))

    return options.length
      ? { from: pos - word.length - 1, to: pos, options: ranked(options) }
      : null
  }

  if (dialect !== 'ql') return null

  const word = /([A-Za-z_$][\w$]*)$/.exec(text.slice(0, pos))?.[1] || ''
  if (!word && !explicit) return null

  const named = new Set(queryFields.map(f => f.name))
  const options: CompletionOption[] = [
    ...queryFields
      .filter(f => startsWith(f.name, word))
      .map(f => ({
        label: f.name,
        detail:
          f.label && f.label !== f.name ? `${f.kind || 'field'} · ${f.label}` : f.kind || 'field',
        boost: BOOST.field,
      })),
    // A module may declare a field of its own with a system name; the module's
    // wins, so the built-in is dropped rather than listed twice.
    ...QL_SYSTEM_FIELDS.filter(n => !named.has(n) && startsWith(n, word)).map(n => ({
      label: n,
      detail: 'record field',
      boost: BOOST.system,
    })),
    ...QL_KEYWORDS.filter(k => startsWith(k, word)).map(k => ({
      label: k,
      detail: 'keyword',
      boost: BOOST.keyword,
    })),
  ]

  return options.length ? { from: pos - word.length, to: pos, options: ranked(options) } : null
}
