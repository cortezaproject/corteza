import { corredor } from '@cortezaproject/corteza-js'

/**
 * Corredor automation mixin
 */

/**
 * With this prefix we distinguish script type
 * @type {string}
 */
const serverScriptPrefix = '/server-scripts/'

export default {
  methods: {
    /**
     * Creates a function for registering server automation scripts to UIHooks and EventBus plugins
     *
     * API should be corteza API Client class that is passed as a first arg to serverScriptHandler
     * See:
     *  - TriggerComposeServerScriptOnManual
     *  - TriggerSystemServerScriptOnManual
     */
    makeAutomationScriptsRegistrator (serverScriptHandler) {
      // eslint-disable-next-line @typescript-eslint/explicit-function-return-type
      return ({ set }) => {
        if (!set || !Array.isArray(set) || set.length === 0) {
          return
        }

        /**
         * Register only server-side scripts (!bundle) and only triggers with onManual eventType
         *
         *  1. client-scripts (bundled) are registered in the bundle's boot loader
         *  2. onManual only -- other kinds (implicit, deferred) are handled directly in/by the Corteza API backend
         */
        set
          .filter(({ name }) => name.substring(0, serverScriptPrefix.length) === serverScriptPrefix)
          .forEach(s => {
            s.triggers
              .filter(({ eventTypes = [] }) => eventTypes.includes('onManual'))
              .forEach(t => {
                // We are (ab)using eventbus for dispatching onManual scripts as well
                // and since it does not know about script structs (only triggers), we need
                // to modify trigger we're passing to it by adding script name
                t.scriptName = s.name
                try {
                  this.$EventBus.Register(ev => serverScriptHandler(ev, s.name), t)
                } catch (e) {
                  console.error(e)
                }
              })
          })

        /**
         * Register all
         */
        this.$UIHooks.Register(...set)
      }
    },

    /**
     * Loads bundle from system API and registers
     * @return {Promise<T>}
     */
    // eslint-disable-next-line @typescript-eslint/explicit-function-return-type
    loadBundle ({ bundle, type = 'client-scripts', ext = 'js', verbose = false, ctx = undefined } = {}) {
      const ep = this.$SystemAPI.automationBundleEndpoint({ bundle, type, ext })

      if (ctx === undefined) {
        throw new Error('can not load bundle and register scripts without context')
      }

      if (typeof ctx.withArgs !== 'function') {
        throw new Error('invalid context object, expecting withArgs function')
      }

      return this.$SystemAPI.api().get(ep)
        .then(({ data }) => {
          if (!data) {
            if (verbose) {
              console.debug('corredor.loadBundle: empty', { bundle, type, ext })
            }
            return
          }

          if (verbose) {
            console.debug('corredor.loadBundle: loaded', { bundle, type, ext })
          }

          // eval loaded bundle
          // eslint-disable-next-line no-new-func
          (new Function(data))()

          if (!window[`${bundle}ClientScripts`]) {
            console.warn(`corredor.loadBundle: window[${bundle}ClientScripts] not defined`)
            return
          }

          const scripts = window[`${bundle}ClientScripts`].scripts || []

          console.debug('corredor.loadBundle:', scripts.length, 'client scripts found')

          scripts.forEach((script) => {
            script.triggers.forEach((trigger) => {
              // Assign script name to handler/trigger:
              // when triggering scripts manually we always trigger a specific script
              trigger.scriptName = script.name
              try {
                this.$EventBus.Register(
                  // Event handler for client-scripts,
                  // convert event arguments and prepare
                  // context for script's execution
                  (ev) => {
                    const args = new corredor.ArgsProxy(ev.args)
                    return corredor.Exec(script, args, ctx.withArgs(args))
                  },
                  trigger,
                )
              } catch (e) {
                console.error(e)
              }
            })
          })
        })
        .catch(({ message }) => {
          console.warn(`could not load client scripts bundle (bundle: ${bundle}, type: ${type}): ${message}`)
        })
    },
  },
}
