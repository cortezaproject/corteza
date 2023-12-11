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
    'import/no-named-default': 'off',
    'new-cap': 'off',
    'vue/component-name-in-template-casing': ['error', 'kebab-case'],
    'vue/no-v-html': 'off',
    'vue/order-in-components': ['error'],
    'comma-dangle': ['error', 'always-multiline'],
    'vue/no-v-html': 'off',
    'no-misleading-character-class': 'off',
    'no-useless-catch': 'off',
    'no-async-promise-executor': 'off',
    'no-case-declarations': 'off',
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
        'CONTENT'
      ],
      alphabetical: false,
    }],
  },
  parserOptions: {
    parser: 'babel-eslint',
  },
}
