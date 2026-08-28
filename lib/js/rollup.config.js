import typescript from 'rollup-plugin-typescript2'
import { nodeResolve } from '@rollup/plugin-node-resolve'
import commonjs from '@rollup/plugin-commonjs'
import json from '@rollup/plugin-json'
import { readFileSync } from 'fs'
import { createRequire } from 'module'

const require = createRequire(import.meta.url)
const pkg = JSON.parse(readFileSync(new URL('./package.json', import.meta.url)))
const ts = require('typescript')

export default {
  input: 'src/index.ts',
  output: [
    {
      file: 'dist/index.cjs',
      format: 'cjs',
      sourcemap: true,
    },
    {
      file: 'dist/index.js',
      format: 'es',
      sourcemap: true,
    },
  ],

  external: [
    ...Object.keys(pkg.dependencies || {}),
    ...Object.keys(pkg.peerDependencies || {}),
    'fs',
    'path',
  ],

  plugins: [
    typescript({
      typescript: ts,
    }),

    nodeResolve({
      browser: true,
      preferBuiltins: true,
      extensions: ['.mjs', '.js', '.json', '.node', '.ts', '.tsx'],
    }),

    commonjs({
      include: /node_modules/,
    }),

    json(),
  ],
}
