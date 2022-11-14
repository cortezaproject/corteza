module.exports = {
  require: [
    'esm',
    'ts-node/register',
    'source-map-support/register',
  ],
  'full-trace': true,
  bail: true,
  recursive: true,
  extension: ['ts', 'js'],
  spec: 'src/**/*.test.ts',
  'watch-files': [ 'src/**' ],
}
