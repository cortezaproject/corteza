/** @type {import('jest').Config} */
module.exports = {
  preset: '@vue/cli-plugin-unit-jest',
  testMatch: [
    '**/tests/unit/**/*.spec.js',
  ],
  moduleNameMapper: {
    '^corteza-webapp-compose/(.*)$': '<rootDir>/$1',
  },
  moduleFileExtensions: [
    'js',
    'jsx',
    'json',
    'vue',
  ],
  transform: {
    '^.+\\.vue$': '@vue/vue2-jest',
    '.+\\.(css|styl|less|sass|scss|svg|png|jpg|ttf|woff|woff2)$': 'jest-transform-stub',
    '^.+\\.jsx?$': 'babel-jest',
  },
  transformIgnorePatterns: [
    '/node_modules/(?!(chai|sinon)/)',
  ],
  snapshotSerializers: [
    'jest-serializer-vue',
  ],
  testEnvironment: 'jsdom',
  testURL: 'http://localhost/',
  watchPlugins: [
    'jest-watch-typeahead/filename',
    'jest-watch-typeahead/testname',
  ],
  setupFiles: [
    '<rootDir>/tests/setup.js',
  ],
}
