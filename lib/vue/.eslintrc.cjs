module.exports = {
  root: true,
  env: {
    node: true,
    es6: true,
    browser: true,
    mocha: true,
  },
  extends: [
    'eslint:recommended',
    'plugin:@typescript-eslint/recommended',
    'plugin:vue/recommended',
  ],
  rules: {
    'no-console': process.env.NODE_ENV === 'production' ? 'error' : 'off',
    'no-debugger': process.env.NODE_ENV === 'production' ? 'error' : 'off',
    'import/no-named-default': 'off',
    'new-cap': 'off',
    'prefer-const': 'off',
    'comma-dangle': ['error', 'always-multiline'],
    'no-misleading-character-class': 'off',
    'no-useless-catch': 'off',
    'no-async-promise-executor': 'off',
    'no-case-declarations': 'off',
    '@typescript-eslint/no-explicit-any': 'off',
    '@typescript-eslint/no-unused-vars': 'off',
    '@typescript-eslint/ban-ts-ignore': 'off',
    '@typescript-eslint/ban-ts-comment': 'off',
    '@typescript-eslint/camelcase': 'off',
    'vue/no-v-model-argument': 'off',
    'vue/component-name-in-template-casing': ['error', 'kebab-case'],
    'vue/order-in-components': ['error'],
    'vue/no-v-html': 'off',
    'vue/multi-word-component-names': 'off',
    'vue/no-lone-template': 'off',
    'vue/no-mutating-props': 'off',
    'vue/component-definition-name-casing': 'off',
    'vue/no-v-text-v-html-on-component': 'off',
    'vue/attributes-order': ['error', {
      order: [
        'DEFINITION',
        'LIST_RENDERING',
        'CONDITIONALS',
        'RENDER_MODIFIERS',
        'GLOBAL',
        'UNIQUE',
        'TWO_WAY_BINDING',
        'SLOT',
        'OTHER_DIRECTIVES',
        'OTHER_ATTR',
        'EVENTS',
        'CONTENT',
      ],
      alphabetical: false,
    }],
  },
  parser: 'vue-eslint-parser',
  parserOptions: {
    parser: '@typescript-eslint/parser',
    ecmaVersion: 2020,
    sourceType: 'module',
    extraFileExtensions: ['.vue'],
    ecmaFeatures: {
      jsx: true,
    },
  },
  plugins: [
    '@typescript-eslint',
    'vue',
  ],
  settings: {
    'import/parsers': {
      '@typescript-eslint/parser': ['.ts', '.tsx'],
    },
    'import/resolver': {
      typescript: {
        alwaysTryTypes: true,
        project: 'tsconfig.json',
      },
    },
  },
  overrides: [
    {
      files: ['*.test.ts'],
      rules: {
        'no-unused-expressions': 'off',
      },
    },
    {
      files: ['*.ts', '*.tsx'],
      rules: {
        'no-undef': 'off', // TypeScript already checks this
      },
    },
    {
      files: ['*.vue'],
      rules: {
        'indent': 'off', // Handled by vue/script-indent
      },
    },
  ],
}
