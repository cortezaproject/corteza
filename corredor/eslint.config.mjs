import tsPlugin from '@typescript-eslint/eslint-plugin'
import sharedConfig from '../eslint.config.shared.js'

export default [
  ...sharedConfig,

  {
    name: 'corredor/files-to-ignore',
    ignores: ['usr/**', 'certs/**'],
  },

  {
    name: 'corredor/project-specific',
    files: ['src/**/*.{js,mjs,ts}'],
    rules: {
      'no-console': process.env.NODE_ENV === 'production' ? 'error' : 'off',
      'no-debugger': process.env.NODE_ENV === 'production' ? 'error' : 'off',
      'no-unused-vars': 'off',
    },
  },

  {
    name: 'corredor/typescript',
    files: ['src/**/*.ts'],
    plugins: {
      '@typescript-eslint': tsPlugin,
    },
    rules: {
      '@typescript-eslint/no-require-imports': 'off',

      // Stand-in for tsc's noUnusedLocals / noUnusedParameters
      '@typescript-eslint/no-unused-vars': [
        'error',
        { argsIgnorePattern: '^_', varsIgnorePattern: '^_', caughtErrorsIgnorePattern: '^_' },
      ],
    },
  },
]
