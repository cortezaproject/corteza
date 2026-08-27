import vue from '@vitejs/plugin-vue'
import { execSync } from 'child_process'
import { readFileSync } from 'fs'
import { fileURLToPath, URL } from 'node:url'
import { defineConfig, loadEnv } from 'vite'
import vueDevTools from 'vite-plugin-vue-devtools'

function getVersion() {
  try {
    return (
      process.env.BUILD_VERSION ||
      execSync('git describe --always --tags', { encoding: 'utf8' }).trim()
    )
  } catch (error) {
    console.error(error)
    return 'unknown'
  }
}

// Read HumanAPI from public/config.js to derive proxy target
function getServerUrl() {
  try {
    const config = readFileSync('./public/config.js', 'utf8')
    const match = config.match(/window\.HumanAPI\s*=\s*['"]([^'"]+)['"]/)
    if (match) {
      return match[1].replace(/\/api\/?$/, '')
    }
  } catch {
    // fallback
  }
  return ''
}

// Vite watches its own `root` as a directory tree, but reaches lib/js and
// lib/vue through workspace symlinks outside it and watches those files one by
// one. A per-file watch does not survive a git write: the first checkout,
// merge or branch switch delivers its event and drops the watch, and every
// later change to that file is then invisible — the module keeps serving the
// transform it had, for the life of the server. When the frozen module is
// lib/vue's barrel, an export added to it is missing from what the browser
// gets and the app dies at import with a blank page.
//
// Watching the source directories instead is what the app's own `src` already
// gets, and a directory watch survives what a per-file one does not.
function watchLibSources() {
  const dirs = ['../../../lib/js/src', '../../../lib/vue/src'].map(dir =>
    fileURLToPath(new URL(dir, import.meta.url)),
  )

  return {
    name: 'human:watch-lib-sources',
    apply: 'serve',
    configureServer(server) {
      // Braces on purpose: vite calls whatever configureServer returns, and a
      // returned watcher is not a function.
      server.watcher.add(dirs)
    },
  }
}

// Where vite serves the app. strictPort is what makes it honest: without it a
// busy port is answered by silently taking the next one, while .env.e2e,
// playwright, dev_ui_verify and the agent toolkit all keep naming the port that
// was asked for — so every browser check drives whatever else is on it.
function devServer(mode) {
  const env = loadEnv(mode, fileURLToPath(new URL('.', import.meta.url)), '')

  return {
    port: Number(process.env.VITE_PORT || env.VITE_PORT || 5173),
    strictPort: true,
    proxy: {
      '/custom.css': getServerUrl(),
      '/code-snippets.js': getServerUrl(),
    },
  }
}

export default defineConfig(({ mode }) => {
  const isDevelopment = mode === 'development'

  return {
    // Root app always uses absolute base path
    base: '/',

    build: {
      // Use '_assets' instead of default 'assets' to avoid conflict
      // with the Go server's embedded web assets route at /assets
      assetsDir: '_assets',

      rollupOptions: {
        output: {
          // Split the large, stable, eagerly-loaded shared vendors out of the
          // app entry chunk so they cache independently across deploys and the
          // entry shrinks. (Heavy section-only libs — tiptap/echarts/fullcalendar
          // /vue-flow/leaflet — are reached only from lazy route chunks, so we
          // leave Rollup's automatic per-route splitting to keep them lazy.)
          manualChunks(id) {
            if (id.includes('node_modules')) {
              if (/node_modules[\\/](primevue|@primevue|@primeuix|primeicons)/.test(id)) {
                return 'vendor-primevue'
              }
              if (/node_modules[\\/](@vue[\\/]|vue[\\/]|vue-router|pinia|vue-i18n)/.test(id)) {
                return 'vendor-vue'
              }
            }
            return undefined
          },
        },
      },
    },

    server: isDevelopment ? devServer(mode) : {},

    plugins: [vue(), vueDevTools(), ...(isDevelopment ? [watchLibSources()] : [])],
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url)),
      },
    },
    define: {
      VERSION: JSON.stringify(getVersion()),
      BUILD_TIME: JSON.stringify(new Date().toISOString()),
    },
  }
})
