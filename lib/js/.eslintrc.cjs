module.exports = {
  root: false,
  env: {
    node: true,
    es6: true,
    mocha: true,
  },
  extends: [
    'eslint:recommended',
    'plugin:@typescript-eslint/recommended',
  ],
  rules: {
    'no-console': process.env.NODE_ENV === 'production' ? 'error' : 'off',
    'no-debugger': process.env.NODE_ENV === 'production' ? 'error' : 'off',
    'comma-dangle': ['error', 'always-multiline'],
    'no-unused-vars': 'off',
    '@typescript-eslint/no-explicit-any': 'off',
    '@typescript-eslint/no-unused-vars': 'off',
    '@typescript-eslint/ban-ts-ignore': 'off',
    '@typescript-eslint/ban-ts-comment': 'off',
    '@typescript-eslint/camelcase': 'off',
  },
  parser: '@typescript-eslint/parser',
  plugins: [
    '@typescript-eslint/eslint-plugin',
  ],
  settings: {
    'import/parsers': {
      '@typescript-eslint/parser': [
        '.ts',
      ],
    },
    'import/resolver': {
      typescript: {},
    },
  },
}
