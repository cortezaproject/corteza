import { compose, system } from '@planetcrust/human-js'
import {
  ComposeCtx,
  ScriptBusPlugin,
  UIHooksPlugin,
  WebappCtx,
  loadClientScripts,
  registerServerScripts,
  usePageStore,
} from '@planetcrust/human-vue'
import router from '../router'

// Apps a manual script can be offered in. `compose` and `admin` are the Corteza
// webapps this one merged; `unify` is Human's own bundle, for scripts that
// belong to no single section.
const apps = ['compose', 'admin', 'unify']

/**
 * Installs the script bus and the manual-script hooks
 */
export function installCorredor(app) {
  app.use(ScriptBusPlugin)
  app.use(UIHooksPlugin, { apps })
}

/**
 * Registers the Corredor scripts this user may run
 *
 * Server scripts are registered from the automation list and executed through
 * the API; client scripts come as a bundle per app and run in the browser.
 */
export async function setupCorredor(app) {
  const globals = app.config.globalProperties
  const { $Auth, $ComposeAPI, $ScriptBus, $SystemAPI, $UIHooks } = globals

  const triggerComposeScript = compose.TriggerComposeServerScriptOnManual($ComposeAPI)
  const triggerSystemScript = system.TriggerSystemServerScriptOnManual($SystemAPI)

  // Both services are asked: each one's automation list drops the scripts bound
  // to the other's resources, and this webapp is the compose and the admin app
  // at once.
  const lists = await Promise.all([
    $ComposeAPI.automationList({ excludeInvalid: true }),
    $SystemAPI.automationList({ excludeInvalid: true }),
  ])

  const byName = new Map()
  lists.forEach(({ set = [] }) => set.forEach(script => byName.set(script.name, script)))

  registerServerScripts({
    scriptBus: $ScriptBus,
    uiHooks: $UIHooks,
    scripts: [...byName.values()],
    handler: (ev, script) =>
      ev.resourceType.startsWith('system')
        ? triggerSystemScript(ev, script)
        : triggerComposeScript(ev, script),
  })

  const args = { $invoker: $Auth.user, authToken: $Auth.accessToken }
  const apis = { systemAPI: $SystemAPI, composeAPI: $ComposeAPI }

  // Read off the globals when a script asks for it: the toast service is
  // installed after this runs.
  const toast = {
    success: message => globals.$toast?.toastSuccess(message),
    warning: message => globals.$toast?.toastWarning(message),
  }

  const composeCtx = new ComposeCtx(args, {
    ...apis,
    pages: () => usePageStore().set,
    toast,
    routePusher: to => router.push(to),
  })

  const webappCtx = new WebappCtx(args, apis)

  const bundles = { compose: composeCtx, admin: webappCtx, unify: webappCtx }

  await Promise.all(
    Object.entries(bundles).map(([bundle, ctx]) =>
      loadClientScripts({ systemAPI: $SystemAPI, scriptBus: $ScriptBus, bundle, ctx }),
    ),
  )
}
