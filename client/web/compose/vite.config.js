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

// Read CortezaAPI from public/config.js to derive proxy target
function getServerUrl() {
  try {
    const config = readFileSync('./public/config.js', 'utf8')
    const match = config.match(/window\.CortezaAPI\s*=\s*['"]([^'"]+)['"]/)
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
    // Set base URL - similar to publicPath in webpack
    base: isDevelopment ? '/' : '/compose/',

    build: {
      assetsDir: '_assets',
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
    optimizeDeps: {
      include: ['echarts', 'vue-echarts'],
    },
    define: {
      VERSION: JSON.stringify(getVersion()),
      BUILD_TIME: JSON.stringify(new Date().toISOString()),
    },
  }
})
