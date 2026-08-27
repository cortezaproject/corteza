const base = require('./jest.config')
const path = require('path')

module.exports = {
  ...base,
  rootDir: __dirname,
  roots: [__dirname, path.resolve(__dirname, '../c311-tests')],
  moduleDirectories: ['node_modules', path.resolve(__dirname, 'node_modules')],
  moduleNameMapper: {
    ...base.moduleNameMapper,
    '^.+/libs/c311-i18n$': path.resolve(__dirname, '../c311-tests/i18n-test-stub.js'),
    '^axios$': path.resolve(__dirname, '../c311-tests/axios-stub.js'),
  },
  testMatch: ['<rootDir>/../c311-tests/**/*.spec.js'],
  setupFiles: [],
}
