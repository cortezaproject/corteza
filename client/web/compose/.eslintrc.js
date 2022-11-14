module.exports = {
  root: true,
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
    'vue/order-in-components': ['error'],
    'comma-dangle': ['error', 'always-multiline'],
    // https://github.com/vuejs/eslint-plugin-vue/blob/master/docs/rules/order-in-components.md
    'vue/no-v-html': 'off',
  },
  parserOptions: {
    parser: 'babel-eslint',
  },
}
