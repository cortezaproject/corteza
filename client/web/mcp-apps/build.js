// Builds every view in src/views into its own self-contained HTML file,
// named <view>.gen.html where the server embeds it.
import { readdirSync, renameSync } from 'node:fs'
import { build } from 'vite'

const src = new URL('./src/views/', import.meta.url)
const out = new URL('../../../server/compose/agentic/mcpui/', import.meta.url)

const views = readdirSync(src)
  .filter(f => f.endsWith('.html'))
  .map(f => f.replace(/\.html$/, ''))

for (const view of views) {
  process.env.VIEW = view
  await build({ configFile: new URL('./vite.config.js', import.meta.url).pathname })
  renameSync(new URL(`${view}.html`, out), new URL(`${view}.gen.html`, out))
}
