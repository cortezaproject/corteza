const base = require('./jest.config')
const path = require('path')

module.exports = {
  ...base,
  preset: path.resolve(__dirname, 'node_modules/@vue/cli-plugin-unit-jest/presets/default'),
  rootDir: path.resolve(__dirname, '../../..'),
  roots: [__dirname, path.resolve(__dirname, '../c311-tests'), path.resolve(__dirname, '../../../lib/vue/src')],
  moduleDirectories: ['node_modules', path.resolve(__dirname, 'node_modules')],
  moduleNameMapper: {
    '^corteza-webapp-compose/(.*)$': '<rootDir>/client/web/compose/$1',
    '^.+/libs/c311-i18n$': path.resolve(__dirname, '../c311-tests/i18n-test-stub.js'),
    '^axios$': path.resolve(__dirname, '../c311-tests/axios-stub.js'),
  },
  transform: {
    '^.+\\.vue$': require.resolve('@vue/vue2-jest'),
    '.+\\.(css|styl|less|sass|scss|svg|png|jpg|ttf|woff|woff2)$': require.resolve('jest-transform-stub'),
    '^.+\\.jsx?$': require.resolve('babel-jest'),
  },
  snapshotSerializers: [require.resolve('jest-serializer-vue')],
  watchPlugins: [
    require.resolve('jest-watch-typeahead/filename'),
    require.resolve('jest-watch-typeahead/testname'),
  ],
  testMatch: ['<rootDir>/client/web/c311-tests/**/*.spec.js'],
  setupFiles: [],
}
