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
  enforced: [
    'client/web/unify/src/App.vue',
    'client/web/unify/src/main.js',
    'client/web/unify/src/config-check.js',
    'client/web/unify/src/router',
    'client/web/unify/src/plugins',
    'client/web/unify/src/utils',
    'client/web/unify/src/sections/index.js',
    'client/web/unify/src/sections/home',
    'client/web/unify/src/sections/admin',
  ],

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
  // Route-target views (view = one screen = one UX contract) get their own
  // sidecar; non-route sub-components stay governed by the resource folder doc.
  // The filename set below matches exactly the components referenced in routes.
  fileTier: [
    'client/web/unify/src/App.vue',
    '**/views/**/List.vue',
    '**/views/**/Editor.vue',
    '**/views/**/Index.vue',
    '**/views/**/View.vue',
    '**/views/**/Configure.vue',
    '**/views/**/Dashboard.vue',
    '**/views/**/Hit.vue',
    '**/views/**/Route.vue',
    '**/views/**/Home.vue',
    '**/views/Dashboard.vue',
    '**/views/Home.vue',
    '**/stores/**/*.js',
    '**/store/**/*.js',
    '**/registry.ts',
  ],

  // Hard cap on prose lines (frontmatter excluded, blank lines not counted).
  docMaxLines: 60,
}
