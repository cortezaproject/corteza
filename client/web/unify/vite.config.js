import vue from '@vitejs/plugin-vue'
import { execSync } from 'child_process'
import { readFileSync } from 'fs'
import { fileURLToPath, URL } from 'node:url'
import { defineConfig } from 'vite'
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

    server: isDevelopment
      ? {
          proxy: {
            '/custom.css': getServerUrl(),
            '/code-snippets.js': getServerUrl(),
          },
        }
      : {},

    plugins: [vue(), vueDevTools()],
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
