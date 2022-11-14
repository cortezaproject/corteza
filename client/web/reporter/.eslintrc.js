module.exports = {
  root: false,
  env: {
    es6: true,
    node: true,
    mocha: true,
  },
  extends: [
    'plugin:vue/recommended',
    '@vue/standard',
  ],
  rules: {
    'no-console': process.env.NODE_ENV === 'production' ? 'error' : 'off',
    'no-debugger': process.env.NODE_ENV === 'production' ? 'error' : 'off',
    'new-cap': 'off',
    'import/no-named-default': 'off',
    'vue/component-name-in-template-casing': ['error', 'kebab-case'],
    'vue/no-v-html': 'off',
    'vue/order-in-components': ['error'],
    'comma-dangle': ['error', 'always-multiline'],
  },
  parserOptions: {
    parser: 'babel-eslint',
  },
}
