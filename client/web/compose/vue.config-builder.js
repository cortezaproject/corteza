var webpack = require('webpack')
var exec = require('child_process').execSync
var path = require('path')
var Vue = require('vue')

module.exports = ({ appFlavour, appLabel, version = process.env.BUILD_VERSION, theme, packageAlias, root = path.resolve('.'), env = process.env.NODE_ENV }) => {
  const isDevelopment = (env === 'development')
  const isTest = (env === 'test')

  if (isTest) {
    Vue.config.devtools = false
    Vue.config.productionTip = false
  }

  if (isDevelopment) {
    Vue.config.devtools = true
    Vue.config.performance = false
  }

  const optimization = isTest ? {} : {
    usedExports: true,
    runtimeChunk: 'single',
    splitChunks: {
      cacheGroups: {
        vendor: {
          test: /[\\/]node_modules[\\/]/,
          chunks: 'all',
        },
      },
    },
  }

  return {
    publicPath: './',
    lintOnSave: true,
    runtimeCompiler: true,

    configureWebpack: {
      // other webpack options to merge in ...
      plugins: [
        new webpack.DefinePlugin({
          FLAVOUR: JSON.stringify(appFlavour),
          WEBAPP: JSON.stringify(appLabel),
          VERSION: JSON.stringify(version || ('' + exec('git describe --always --tags')).trim()),
          BUILD_TIME: JSON.stringify((new Date()).toISOString()),
        }),
      ],

      optimization,
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
            const regexp = /\[mini-css-extract-plugin\]/
            if (regexp.test(error.message)) return {}
            return error
          },
        ]
        return args
      })

      // Do not copy config files (deployment procedure will do that)
      config.plugins.has('copy') && config.plugin('copy').tap(options => {
        options[0][0].ignore.push('config*js')
        return options
      })

      // Aliasing 'corteza-webapp-compose' instead of '@' so we do
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
        }))

      // Load CSS assets according to their location
      scssNormal.use('resolve-url-loader')
        .loader('resolve-url-loader').options({
          keepQuery: true,
          removeCR: true,
          root: path.join(root, 'src/themes', theme),
        })
        .before('sass-loader')
    },

    devServer: {
      host: '127.0.0.1',
      hot: true,
      disableHostCheck: true,

      watchOptions: {
        ignored: [
          // Do not watch for changes under node_modules
          // (exception is node_modules/@cortezaproject)
          /node_modules([\\]+|\/)+(?!@cortezaproject)/,
        ],
        aggregateTimeout: 200,
        poll: 1000,
      },
    },

    css: {
      sourceMap: isDevelopment,
      loaderOptions: {
        sass: {
          // @todo cleanup all components and remove this global import
          additionalData: `@import "./src/themes/${theme}/variables.scss";`,
        },
      },
    },
  }
}
