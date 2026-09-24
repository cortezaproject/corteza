import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'
import { defineConfig } from 'vite'
import { viteSingleFile } from 'vite-plugin-singlefile'

// One self-contained HTML file per view, written where the Go server embeds it.
// VIEW picks the view: `VIEW=record_lookup vite build`.
const view = process.env.VIEW || 'record_lookup'

export default defineConfig({
  root: fileURLToPath(new URL('./src/views', import.meta.url)),
  plugins: [vue(), viteSingleFile()],
  build: {
    outDir: fileURLToPath(new URL('../../../server/compose/agentic/mcpui', import.meta.url)),
    emptyOutDir: false,
    rollupOptions: {
      input: fileURLToPath(new URL(`./src/views/${view}.html`, import.meta.url)),
    },
  },
})
