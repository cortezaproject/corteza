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

  // Shared rules for all projects
  {
    name: 'shared/rules',
    rules: {
      // Indentation and formatting
      indent: ['error', 2, { SwitchCase: 1 }],
      'vue/html-indent': ['error', 2],
      'vue/script-indent': ['error', 2, { baseIndent: 0 }],

      // Quote enforcement
      quotes: ['error', 'single', { avoidEscape: true }],
      'vue/html-quotes': ['error', 'double', { avoidEscape: false }],

      // Vue specific rules
      'vue/max-attributes-per-line': [
        'error',
        {
          singleline: { max: 3 },
          multiline: { max: 1 },
        },
      ],
      'vue/multiline-html-element-content-newline': 'error',
      'vue/singleline-html-element-content-newline': 'off',
      'vue/multi-word-component-names': 'off',
      'vue/no-reserved-component-names': 'off',

      // General JavaScript rules
      'no-unused-vars': 'warn',
      'no-console': 'off',
      semi: ['error', 'never'],
      'comma-dangle': ['error', 'always-multiline'],
      'object-curly-spacing': ['error', 'always'],
      'array-bracket-spacing': ['error', 'never'],
      'computed-property-spacing': ['error', 'never'],
      'key-spacing': ['error', { beforeColon: false, afterColon: true }],
      'space-before-function-paren': ['error', 'never'],
    },
  },
]
