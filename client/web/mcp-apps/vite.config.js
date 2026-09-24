import vue from '@vitejs/plugin-vue'
import { load } from 'js-yaml'
import { readFileSync } from 'node:fs'
import { fileURLToPath, URL } from 'node:url'
import { defineConfig } from 'vite'
import { viteSingleFile } from 'vite-plugin-singlefile'

// One self-contained HTML file per view, written where the Go server embeds it.
// VIEW picks the view: `VIEW=record_lookup vite build`.
const view = process.env.VIEW || 'record_lookup'

// The human-webapp translation files a view reads, keyed by file name the way
// the server's merged bundle keys them. A view has no API to fetch them from.
const localeFiles = ['field', 'general', 'mcpApp']

function humanLocale() {
  const id = 'virtual:human-locale'
  return {
    name: 'human:locale',
    resolveId: source => (source === id ? '\0' + id : undefined),
    load(source) {
      if (source !== '\0' + id) return
      const en = {}
      for (const name of localeFiles) {
        const file = new URL(`../../../locale/en/human-webapp/${name}.yaml`, import.meta.url)
        this.addWatchFile(fileURLToPath(file))
        en[name] = load(readFileSync(file, 'utf8'))
      }
      return `export default ${JSON.stringify({ en })}`
    },
  }
}

export default defineConfig({
  root: fileURLToPath(new URL('./src/views', import.meta.url)),
  plugins: [vue(), humanLocale(), viteSingleFile()],
  build: {
    outDir: fileURLToPath(new URL('../../../server/compose/agentic/mcpui', import.meta.url)),
    emptyOutDir: false,
    rollupOptions: {
      input: fileURLToPath(new URL(`./src/views/${view}.html`, import.meta.url)),
    },
  },
})
