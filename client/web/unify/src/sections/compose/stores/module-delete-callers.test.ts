import { describe, expect, it } from 'vitest'
import { readFileSync, readdirSync, statSync } from 'node:fs'
import { join, resolve } from 'node:path'

// The API client refuses moduleDelete without a namespaceID, throwing
// "field namespaceID is empty" before the request is ever made. That guard is
// good — but it fires at run time, in a toast, and the store happily forwards
// whatever a caller hands it.
//
// Deleting a module from the edit view did exactly that for as long as the view
// existed: it passed only moduleID, so the action failed for every user while
// the list view, which passes both, worked fine (issue #30).
//
// A unit test of the store would not have caught it — module.test.ts passes
// both fields and always did. The gap was the caller, and no view in this
// section is mounted by any test, so this checks the call sites at the source
// level instead. It is a blunt instrument, and it is the one that would have
// caught the bug.
const roots = [
  resolve(__dirname, '../../..'),
  resolve(__dirname, '../../../../../../../lib/vue/src'),
]

const callPattern = /(?:moduleStore\.delete|\$?ComposeAPI\.moduleDelete|api\.moduleDelete)\s*\(/g

function sourceFiles(dir: string): string[] {
  const out: string[] = []

  for (const entry of readdirSync(dir)) {
    if (entry === 'node_modules' || entry === 'dist' || entry.startsWith('.')) {
      continue
    }

    const path = join(dir, entry)

    if (statSync(path).isDirectory()) {
      out.push(...sourceFiles(path))
      continue
    }

    if (/\.(vue|ts|js)$/.test(entry) && !/\.test\.ts$/.test(entry)) {
      out.push(path)
    }
  }

  return out
}

// argumentsOf returns the text between the call's parentheses, so a multi-line
// argument object is read whole rather than by line.
function argumentsOf(source: string, openIndex: number): string {
  let depth = 0

  for (let i = openIndex; i < source.length; i++) {
    if (source[i] === '(') depth++
    if (source[i] === ')') {
      depth--
      if (depth === 0) return source.slice(openIndex + 1, i)
    }
  }

  return source.slice(openIndex)
}

describe('module delete callers', () => {
  it('every call site passes a namespaceID', () => {
    const offenders: string[] = []

    for (const root of roots) {
      for (const file of sourceFiles(root)) {
        const source = readFileSync(file, 'utf8')

        for (const match of source.matchAll(callPattern)) {
          const open = match.index! + match[0].length - 1
          const args = argumentsOf(source, open)

          // A call forwarding an already-built object (deleteModule(item)) is
          // not a call site that can be checked here — its caller is.
          if (/^\s*[a-zA-Z_$][\w$]*\s*(,|$)/.test(args)) {
            continue
          }

          if (!args.includes('namespaceID')) {
            const line = source.slice(0, match.index!).split('\n').length
            offenders.push(`${file.replace(process.cwd(), '.')}:${line}`)
          }
        }
      }
    }

    expect(offenders, 'moduleDelete requires namespaceID; these call sites omit it').toEqual([])
  })
})
