import vue from 'rollup-plugin-vue'
import resolve from '@rollup/plugin-node-resolve'
import commonjs from '@rollup/plugin-commonjs'
import typescript from 'rollup-plugin-typescript2'
import babel from '@rollup/plugin-babel'
import json from '@rollup/plugin-json'
import styles from 'rollup-plugin-styles'

import { readFileSync } from 'fs'
import { createRequire } from 'module'

const require = createRequire(import.meta.url)
const pkg = JSON.parse(readFileSync(new URL('./package.json', import.meta.url)))
const ts = require('typescript')

export default {
  input: 'src/index.ts',
  output: [
    {
      file: pkg.main,
      format: 'cjs',
      sourcemap: true,
      inlineDynamicImports: true,
    },
    {
      file: pkg.module,
      format: 'es',
      sourcemap: true,
      inlineDynamicImports: true,
    },
  ],

  external: [
    ...Object.keys(pkg.dependencies || {}),
    ...Object.keys(pkg.peerDependencies || {}),
    'fs',
    'path',
  ],

  plugins: [
    resolve({
      browser: true,
      preferBuiltins: false,
    }),
    commonjs({
      include: /node_modules/,
    }),
    typescript({
      typescript: ts,
      tsconfig: './tsconfig.json',
      sourceMap: true,
      check: true,
    }),
    vue({
      preprocessStyles: true,
    }),
    babel({
      exclude: /node_modules\/(?!pdfjs-dist).*/,
      babelHelpers: 'bundled',
      presets: [
        ['@babel/preset-env'],
      ],
      plugins: [
        '@babel/plugin-transform-class-properties',
        '@babel/plugin-transform-private-methods',
        '@babel/plugin-transform-private-property-in-object',
      ],
    }),
    json(),
    styles(),
  ],

  watch: {
    exclude: ['node_modules/**'],
  },
}
