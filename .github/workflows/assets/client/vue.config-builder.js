const webpack = require('webpack')
const exec = require('child_process').execSync
const path = require('path')
const Vue = require('vue')

module.exports = ({ appFlavour, appLabel, version = process.env.BUILD_VERSION, theme, packageAlias, root = path.resolve('.'), env = process.env.NODE_ENV }) => {
  const isDevelopment = (env === 'development')
  const isTest = (env === 'test')

  if (isTest) {
    Vue.config.devtools = false
    Vue.config.productionTip = false
  }

  if (isDevelopment) {
    Vue.config.devtools = true
    Vue.config.performance = true
  }

  const optimization = isTest
    ? {}
    : {
        usedExports: true,
        runtimeChunk: 'single',
        splitChunks: {
          chunks: 'all',
          minSize: 20000,
          maxSize: 244000,
          cacheGroups: {
            vendor: {
              test: /[\\/]node_modules[\\/]/,
              name: 'vendors',
              chunks: 'all',
              priority: -10,
            },
            common: {
              minChunks: 2,
              priority: -20,
              reuseExistingChunk: true,
            },
          },
        },
      }

  return {
    publicPath: isDevelopment ? '/' : './',
    lintOnSave: true,
    runtimeCompiler: true,

    configureWebpack: {
      // other webpack options to merge in ...

      // Webpack 5 specific configuration
      resolve: {
        symlinks: false,
        fallback: {
          path: false,
          fs: false, // Add explicit false for Node.js core modules
          crypto: false,
          stream: false,
          util: false,
        },
      },

      plugins: [
        new webpack.DefinePlugin({
          FLAVOUR: JSON.stringify(appFlavour),
          WEBAPP: JSON.stringify(appLabel),
          VERSION: JSON.stringify(version || ('' + exec('git describe --always --tags')).trim()),
          BUILD_TIME: JSON.stringify((new Date()).toISOString()),
        }),
      ],

      optimization: {
        ...optimization,
        moduleIds: 'deterministic',
        chunkIds: 'named',
        emitOnErrors: false,
      },

      stats: {
        warnings: false,
        errors: false,
      },

      performance: {
        hints: isDevelopment ? false : 'warning',
        maxEntrypointSize: 512000,
        maxAssetSize: 512000,
      },
    },

    chainWebpack: config => {
      // https://cli.vuejs.org/guide/troubleshooting.html#symbolic-links-in-node-modules
      config.resolve.symlinks(false)

      // Remove css extraction issues
      // https://github.com/vuejs/vue-cli/issues/3771#issuecomment-526228100
      config.plugin('friendly-errors').tap(args => {
        const vueCli3Transformer = args[0].additionalTransformers[0]
        args[0].additionalTransformers = [
          vueCli3Transformer,
          error => {
            return /\[mini-css-extract-plugin\]/.test(error.message) ? {} : error
          },
        ]
        return args
      })

      // Update copy-webpack-plugin configuration for v8+
      config.plugins.has('copy') && config.plugin('copy').tap(options => {
        options[0].patterns = [{
          from: 'public',
          globOptions: {
            ignore: ['**/config*js', '**/*gitignore', '**/index.html'],
          },
          noErrorOnMissing: true,
        }]
        return options
      })

      // Aliasing full package name instead of '@' so we do
      // not break imports on apps that import this code
      config.resolve.alias.delete('@')
      if (packageAlias) {
        config.resolve.alias.set(packageAlias, root)
      }

      if (isTest) {
        const scssRule = config.module.rule('scss')
        scssRule.uses.clear()
        scssRule
          .use('null-loader')
          .loader('null-loader')
      }

      const scssNormal = config.module.rule('scss').oneOf('normal')

      scssNormal.use('sass-loader')
        .loader('sass-loader')
        .tap(options => ({
          ...options,
          sourceMap: true,
          sassOptions: {
            outputStyle: isDevelopment ? 'expanded' : 'compressed',
          },
        }))

      // Load CSS assets according to their location
      scssNormal.use('resolve-url-loader')
        .loader('resolve-url-loader').options({
          keepQuery: true,
          removeCR: true,
          root: path.join(root, 'src/themes', theme),
        })
        .before('sass-loader')

      // Keep this to ensure we don't have multiple HTML plugins
    },

    devServer: {
      host: '127.0.0.1',
      port: 8080,
      hot: true,
      allowedHosts: 'all', // replaces disableHostCheck
      webSocketServer: 'ws',
      headers: {
        'Access-Control-Allow-Origin': '*',
      },

      proxy: {
        '^/custom.css': {
          target: fetchBaseUrl(),
        },

        '^/code-snippets.js': {
          target: fetchBaseUrl(),
        },
      },

      // Webpack 5 DevServer configuration
      watchFiles: {
        paths: [
          '**/*',
          '!**/node_modules/!(@cortezaproject)/**',
        ],
        options: {
          usePolling: true,
          aggregateTimeout: 200,
          poll: 1000,
        },
      },

      client: {
        overlay: false,
        progress: true,
      },
    },

    css: {
      sourceMap: isDevelopment,
      extract: !isTest,
      loaderOptions: {
        sass: {
          sassOptions: {
            outputStyle: isDevelopment ? 'expanded' : 'compressed',
          },
        },
        postcss: {
          postcssOptions: {
            plugins: [
              require('autoprefixer'),
              require('rtlcss'),
            ],
          },
        },
      },
    },
  }
}

function fetchBaseUrl () {
  const fs = require('fs')
  const window = {}

  const fileContents =
    fs.existsSync('public/config.js')
      ? fs.readFileSync('public/config.js', 'utf-8')
      : ''

  try {
    // eslint-disable-next-line no-eval
    eval(fileContents)

    const u = window.CortezaAPI || ''
    const ur = new URL(u.startsWith('//') ? `http:${u}` : u)

    return `${ur.protocol}//${ur.host}/`
  } catch (e) {
    return '/'
  }
}
