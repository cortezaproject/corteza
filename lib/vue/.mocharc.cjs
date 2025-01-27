module.exports = {
  require: [
    'tsx/cjs'
  ],
  'full-trace': true,
  bail: true,
  recursive: true,
  extension: ['ts'],
  spec: ['src/**/*.test.ts'],
  timeout: 5000,
}
