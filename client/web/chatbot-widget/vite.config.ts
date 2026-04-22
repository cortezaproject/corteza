import { defineConfig } from 'vite'
import { fileURLToPath, URL } from 'node:url'

export default defineConfig({
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },

  build: {
    target: 'es2019',
    minify: 'esbuild',
    cssCodeSplit: false,
    emptyOutDir: true,
    outDir: 'dist',
    rollupOptions: {
      input: fileURLToPath(new URL('./src/entry.ts', import.meta.url)),
      output: {
        format: 'iife',
        entryFileNames: 'widget.js',
        inlineDynamicImports: true,
        // CSS is read as a string and injected into shadow root at runtime
        // (see styles.ts), no separate file emitted.
      },
    },
  },

  server: {
    port: 5190,
  },
})
