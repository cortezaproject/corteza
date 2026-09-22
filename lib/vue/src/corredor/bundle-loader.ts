import { corredor } from '@planetcrust/human-js'
import type { Ctx } from '@planetcrust/human-js/src/corredor/ctx'
import type { BaseArgs } from '@planetcrust/human-js/src/corredor/shared'
import type { Script, ScriptTrigger, UIHooks } from './ui-hooks'

const serverScriptPrefix = '/server-scripts/'
const onManual = 'onManual'

export interface ScriptEvent {
  resourceType: string
  eventType: string
  args?: { [_: string]: unknown }
}

export interface ScriptBus {
  Register(_handler: (_ev: ScriptEvent) => Promise<unknown>, _trigger: ScriptTrigger): unknown
}

export interface BundleAPI {
  automationBundleEndpoint(_a: { bundle: string; type: string; ext: string }): string
  api(): { get(_url: string): Promise<{ data?: unknown }> }
}

export interface ScriptCtx {
  withArgs(_args: BaseArgs): Ctx
}

export interface LoadClientScriptsOptions {
  systemAPI: BundleAPI
  scriptBus: ScriptBus
  bundle: string
  ctx: ScriptCtx
  type?: string
  ext?: string
  verbose?: boolean
}

interface ClientScriptBundle {
  scripts?: Array<Script & { exec?: unknown }>
}

/**
 * Loads a client-script bundle and registers every script's triggers on the bus
 *
 * Resolves with the scripts it registered; a bundle that is missing, empty or
 * broken is reported and resolves empty.
 */
export async function loadClientScripts({
  systemAPI,
  scriptBus,
  bundle,
  ctx,
  type = 'client-scripts',
  ext = 'js',
  verbose = false,
}: LoadClientScriptsOptions): Promise<Script[]> {
  const endpoint = systemAPI.automationBundleEndpoint({ bundle, type, ext })

  try {
    const { data } = await systemAPI.api().get(endpoint)

    if (!data || typeof data !== 'string') {
      if (verbose) console.debug('corredor: empty client script bundle', { bundle, type, ext })
      return []
    }

    // The bundle is a webpack build that parks its exports on the global object.
    new Function(data)()

    const globals = window as unknown as { [_: string]: ClientScriptBundle | undefined }
    const loaded = globals[`${bundle}ClientScripts`]

    if (!loaded) {
      console.warn(`corredor: window[${bundle}ClientScripts] not defined`)
      return []
    }

    const scripts = loaded.scripts || []

    if (verbose) console.debug(`corredor: ${scripts.length} client scripts in ${bundle} bundle`)

    scripts.forEach(script => {
      ;(script.triggers || []).forEach(trigger => {
        // Triggering a script manually always triggers one specific script.
        trigger.scriptName = script.name

        try {
          scriptBus.Register(ev => {
            const args = new corredor.ArgsProxy(ev.args || {}) as unknown as BaseArgs

            return corredor.Exec(script as corredor.ExecutableScript, args, ctx.withArgs(args))
          }, trigger)
        } catch (e) {
          console.warn(`corredor: could not register trigger of ${script.name}`, e)
        }
      })
    })

    return scripts
  } catch (e) {
    const { message } = e as Error
    console.warn(`corredor: could not load client script bundle ${bundle} (${type}): ${message}`)
    return []
  }
}

export interface ServerScriptHandler {
  (_ev: ScriptEvent, _script: string): Promise<unknown>
}

export interface RegisterServerScriptsOptions {
  scriptBus: ScriptBus
  uiHooks: UIHooks
  scripts: Script[]
  handler: ServerScriptHandler
}

/**
 * Registers the manually triggered server scripts on the bus and all scripts as
 * UI hooks
 *
 * Only server scripts are registered on the bus: client scripts register
 * themselves when their bundle loads, and implicit and deferred server scripts
 * are handled by the API. The UI hooks take every script, since a button can
 * trigger either kind.
 */
export function registerServerScripts({
  scriptBus,
  uiHooks,
  scripts,
  handler,
}: RegisterServerScriptsOptions): void {
  if (!Array.isArray(scripts) || scripts.length === 0) {
    return
  }

  scripts
    .filter(({ name }) => name.startsWith(serverScriptPrefix))
    .forEach(s => {
      ;(s.triggers || [])
        .filter(({ eventTypes }) => eventTypes?.includes(onManual))
        .forEach(trigger => {
          trigger.scriptName = s.name

          try {
            scriptBus.Register(ev => handler(ev, s.name), trigger)
          } catch (e) {
            console.warn(`corredor: could not register trigger of ${s.name}`, e)
          }
        })
    })

  uiHooks.Register(...scripts)
}
