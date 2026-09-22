module.exports = {
  require: ['tsx/cjs'],
  'full-trace': true,
  bail: true,
  recursive: true,
  extension: ['ts', 'js'],
  spec: 'src/**/*.test.ts',
  'watch-files': ['src/**'],
}
