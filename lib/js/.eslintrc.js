module.exports = {
  root: false,
  env: {
    node: true,
    es6: true,
    mocha: true,
  },
  extends: [
    'standard',
    'plugin:@typescript-eslint/recommended',
  ],
  rules: {
    'no-console': process.env.NODE_ENV === 'production' ? 'error' : 'off',
    'no-debugger': process.env.NODE_ENV === 'production' ? 'error' : 'off',
    'comma-dangle': [ 'error', 'always-multiline' ],
  },
  parser: '@typescript-eslint/parser',
  plugins: [
    '@typescript-eslint',
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
  overrides: [
    {
      "files": ["*.test.ts"],
      "rules": {
          "no-unused-expressions": "off"
      }
    }
  ]
}

