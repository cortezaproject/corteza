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
    'client/web/unify',
    'client/web/chatbot-widget',
    'lib/vue',
    'lib/js',
    'lib/test-utils',
  ],

  // Never covered, never checked (glob-ish: * = segment, ** = any depth).
  exclude: [
    '**/node_modules/**',
    '**/vendor/**',
    '**/dist/**',
    '**/.build/**',
    '**/*.gen.*',
    '**/assets/**',
    'client/web/unify/public/**', // runtime instance config + static assets
    'client/web/unify/e2e/.results/**', // playwright run artifacts (gitignored)
    'client/web/unify/e2e/.auth/**', // playwright session state (gitignored)
    'client/web/chatbot-widget/public/**',
    'lib/js/src/api-clients/**', // codegen-owned (generated from server rest.yaml)
    'extra/**',

    // Project-feature code lives inside shared server packages rather than in a
    // package of its own, so it is excluded by filename. Ruled 2026-08-04; see
    // TODO.md Phase 4. Note `project*.go`, not `project_*.go` — the snake_case
    // form misses system/rest/project.go, system/service/projectTask.go and
    // nine other hand-written files.
    'server/**/project*.go',
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
    '**/views/**/Edit.vue',
    '**/views/**/Create.vue',
    '**/views/**/Builder.vue',
    '**/views/**/RecordView.vue',
    '**/views/**/Sessions.vue',
    '**/views/Dashboard.vue',
    '**/views/Home.vue',
    '**/views/ProjectList.vue',
    '**/views/Wizard.vue',
    '**/views/dashboard/*.vue',
    '**/stores/**/*.js',
    '**/stores/**/*.ts',
    '**/store/**/*.js',
    '**/registry.ts',
  ],

  // Matched against fileTier hits to exempt them from the sidecar requirement.
  fileTierExclude: [
    '**/*.test.*',
    'lib/test-utils/**', // fixture/mock helpers, not real stores/views
  ],

  // Hard cap on prose lines (frontmatter excluded, blank lines not counted).
  // Raised from 60 (2026-08-04): the longest doc in the repo was 90 lines and
  // several were fighting the cap. 250 is ~3x that — it stops constraining a
  // genuinely complex package while still failing a runaway one.
  docMaxLines: 250,
}
