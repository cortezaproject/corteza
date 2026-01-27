import skipFormatting from '@vue/eslint-config-prettier/skip-formatting'
import sharedConfig from '../../eslint.config.shared.js'

export default [
  ...sharedConfig,

  // Project-specific overrides for lib/vue (JS and Vue files only, not TS)
  {
    name: 'lib-vue/project-specific',
    files: ['src/**/*.{js,mjs,jsx,vue}'],
    rules: {
      // Lib-specific rules (stricter for library code)
      'no-unused-vars': ['warn', { argsIgnorePattern: '^_', varsIgnorePattern: '^_' }],
      'no-console': 'warn', // Libraries shouldn't have console logs
    },
  },

  skipFormatting,
]
