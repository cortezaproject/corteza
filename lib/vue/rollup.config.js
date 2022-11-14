import vue from 'rollup-plugin-vue'
import resolve from '@rollup/plugin-node-resolve'
import commonjs from '@rollup/plugin-commonjs'
import typescript from 'rollup-plugin-typescript2'
import babel from 'rollup-plugin-babel'
import json from 'rollup-plugin-json'
import styles from "rollup-plugin-styles"

import pkg from './package.json'

export default {
  input: 'src/index.ts',
  output: [
    {
      file: pkg.main,
      format: 'cjs',
      sourcemap: true,
    },
    {
      file: pkg.module,
      format: 'es',
      sourcemap: true,
    },
  ],

  external: [
    ...Object.keys(pkg.dependencies || {}),
    ...Object.keys(pkg.peerDependencies || {}),
  ],

  plugins: [
    resolve({
      main: true,
      browser: true,
      preferBuiltins: true,
    }),
    typescript({
      typescript: require('typescript'),

      // Preventing :
      // [!] (plugin rpt2) Error: Unknown object type "asyncfunction"
      //
      // see: https://github.com/ezolenko/rollup-plugin-typescript2/issues/105
      objectHashIgnoreUnknownHack: true,
    }),
    commonjs(),
    vue({
      compileTemplate: true,
    }),
    babel({
      exclude: 'node_modules/**',
    }),
    json(),
    styles(),
  ],

  watch: {
    exclude: ['node_modules/**'],
  },
}
