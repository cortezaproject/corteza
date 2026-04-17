import skipFormatting from '@vue/eslint-config-prettier/skip-formatting'
import js from '@eslint/js'
import tsPlugin from '@typescript-eslint/eslint-plugin'
import tsParser from '@typescript-eslint/parser'
import pluginVue from 'eslint-plugin-vue'
import globals from 'globals'

export default [
  {
    name: 'shared/files-to-lint',
    files: ['**/*.{js,mjs,jsx,ts,tsx,vue}'],
  },

  {
    name: 'shared/files-to-ignore',
    ignores: ['**/dist/**', '**/dist-ssr/**', '**/coverage/**', '**/node_modules/**'],
  },

  {
    languageOptions: {
      globals: {
        ...globals.browser,
        ...globals.node,
        ...globals.es2022,
      },
    },
  },

  js.configs.recommended,
  ...pluginVue.configs['flat/essential'],

  // Let vue-eslint-parser route <script lang="ts"> blocks through tsParser
  {
    name: 'shared/vue-ts-parser',
    files: ['**/*.vue'],
    languageOptions: {
      parserOptions: {
        parser: tsParser,
        ecmaVersion: 2022,
        sourceType: 'module',
      },
    },
  },

  // TypeScript configuration
  {
    name: 'shared/typescript',
    files: ['**/*.ts', '**/*.tsx'],
    languageOptions: {
      parser: tsParser,
      parserOptions: {
        ecmaVersion: 2022,
        sourceType: 'module',
      },
    },
    plugins: {
      '@typescript-eslint': tsPlugin,
    },
    rules: {
      ...tsPlugin.configs.recommended.rules,

      // Disable base ESLint rules that conflict with TypeScript
      'no-unused-vars': 'off',
      'no-undef': 'off', // TypeScript handles this

      // TypeScript overrides
      '@typescript-eslint/no-explicit-any': 'off',
      '@typescript-eslint/no-unused-vars': [
        'warn',
        {
          argsIgnorePattern: '^_',
          varsIgnorePattern: '^_',
        },
      ],
      '@typescript-eslint/ban-ts-comment': 'off',
    },
  },

  // Test files configuration
  {
    name: 'shared/test-files',
    files: ['**/*.test.{js,ts,tsx,vue}', '**/*.spec.{js,ts,tsx,vue}'],
    languageOptions: {
      globals: {
        ...globals.mocha,
      },
    },
    rules: {
      // Relax rules for test files
      '@typescript-eslint/no-unused-expressions': 'off',
      'no-unused-expressions': 'off',
      'no-constant-binary-expression': 'off',
    },
  },

  // Shared rules for all projects (code quality only - Prettier handles formatting)
  {
    name: 'shared/rules',
    rules: {
      // Vue specific rules (structural, not formatting)
      'vue/multi-word-component-names': 'off',
      'vue/no-reserved-component-names': 'off',

      // Code quality rules
      'no-unused-vars': ['warn', { argsIgnorePattern: '^_', varsIgnorePattern: '^_', caughtErrorsIgnorePattern: '^_' }],
      'no-console': 'off',
    },
  },

  // Disable all formatting rules - Prettier handles formatting
  skipFormatting,
]
