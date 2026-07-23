// Intent system coverage manifest — see .intent/SPEC.md
export default {
  // Roots that will eventually carry intent docs (backfill target).
  covered: [
    'client/web/unify',
    'client/web/chatbot-widget',
    'lib/vue',
    'lib/js',
    'lib/test-utils',
    'lib/eslint-client',
    'server',
    'def',
    'tests',
    'locale',
  ],

  // Roots (dirs or single files) where check/sync/coverage are ENFORCED today.
  // Grows as backfill phases land. Phase 1 = friction pilot.
  enforced: [],

  // Never covered, never checked (glob-ish: * = segment, ** = any depth).
  exclude: [
    '**/node_modules/**',
    '**/vendor/**',
    '**/dist/**',
    '**/.build/**',
    '**/*.gen.*',
    '**/assets/**',
    'client/web/unify/src/sections/project/**', // POC, still iterating
    'extra/**',
  ],

  // File extensions that count as covered source.
  sourceExt: ['.vue', '.js', '.mjs', '.cjs', '.ts', '.go', '.cue', '.scss', '.css', '.html'],

  // Files that must have their own <name>.intent.md sidecar (beyond folder docs).
  // Globs, matched against repo-relative paths.
  // Admin-style CRUD views are governed per-resource-folder instead (List+Editor
  // share one INTENT.md); sidecars are for singular load-bearing files.
  fileTier: [
    'client/web/unify/src/App.vue',
    '**/stores/**/*.js',
    '**/store/**/*.js',
    '**/registry.ts',
  ],

  // Hard cap on prose lines (frontmatter excluded, blank lines not counted).
  docMaxLines: 60,
}
