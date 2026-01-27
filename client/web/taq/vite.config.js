import vue from '@vitejs/plugin-vue'
import { execSync } from 'child_process'
import { fileURLToPath, URL } from 'node:url'
import { PrimeVueResolver } from 'unplugin-vue-components/resolvers'
import Components from 'unplugin-vue-components/vite'
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

export default defineConfig(({ mode }) => {
  const isDevelopment = mode === 'development'

  return {
    // Set base URL - similar to publicPath in webpack
    base: isDevelopment ? '/' : './',

    // Handle CommonJS dependencies (dagre uses require())
    optimizeDeps: {
      include: ['dagre'],
    },

    plugins: [
      vue(),
      vueDevTools(),
      Components({
        resolvers: [PrimeVueResolver()],
        // Auto-import components from these directories
        dirs: ['src/components'],
        // Generate TypeScript declaration file
        dts: true,
      }),
    ],
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
