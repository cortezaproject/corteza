import { eventbus } from '@planetcrust/human-js'
import type { App, InjectionKey, Plugin } from 'vue'

const onManual = 'onManual'

interface KV {
  [_: string]: string
}

export interface UIProp {
  name: string
  value: string
}

export interface ScriptTrigger {
  resourceTypes?: string[]
  eventTypes?: string[]
  uiProps?: UIProp[]
  constraints?: object[]
  weight?: number
  scriptName?: string
}

export interface Script {
  name: string
  label?: string
  description?: string
  errors?: string[]
  triggers?: ScriptTrigger[]
}

function prop2map(uiProps?: UIProp[]): KV {
  if (!uiProps) {
    return {}
  }

  return uiProps.reduce((m: KV, { name, value }) => {
    m[name] = value
    return m
  }, {})
}

function sorter(a: Button, b: Button): number {
  if (a.weight === b.weight) {
    return a.script.localeCompare(b.script)
  }

  return a.weight - b.weight
}

/**
 * A manually triggered script, as the UI offers it.
 */
export class Button {
  readonly label: string
  readonly description?: string
  readonly script: string
  readonly resourceType: string
  readonly app: string
  readonly weight: number
  readonly variant?: string
  readonly page?: string
  readonly slot?: string
  readonly constraints: eventbus.ConstraintMatcher[]

  constructor(s: Script, t: ScriptTrigger) {
    const uiProps = prop2map(t.uiProps)

    if (!t.eventTypes?.includes(onManual)) {
      throw new Error('expecting onManual event type')
    }

    if (t.resourceTypes?.length !== 1) {
      throw new Error('expecting exactly one resource type on trigger')
    }

    this.label = uiProps.label ?? s.label ?? s.name
    this.description = s.description
    this.script = s.name
    this.weight = t.weight || 0
    this.resourceType = t.resourceTypes[0]
    this.app = uiProps.app ?? ''
    this.page = uiProps.page
    this.slot = uiProps.slot
    this.variant = uiProps.variant
    this.constraints = (t.constraints || []).map(eventbus.ConstraintMaker)
  }
}

export interface UIHooksOptions {
  apps: string[]
  verbose: boolean
}

/**
 * Consumes scripts that can be triggered manually and converts them to buttons.
 *
 * The buttons are either picked by hand into a compose page block, or placed
 * automatically on the page and slot their trigger names.
 */
export class UIHooks {
  readonly apps: string[]
  readonly verbose: boolean

  protected set: Button[] = []

  constructor(opt?: Partial<UIHooksOptions>) {
    this.apps = opt?.apps || []
    this.verbose = !!opt?.verbose
  }

  /**
   * Takes one or more scripts and converts them to buttons.
   *
   * With every script added it removes ALL buttons that use the same script.
   *
   * A trigger naming an app this instance does not serve is skipped; one naming
   * no app at all is kept, so a script can be reached by name from anywhere.
   */
  Register(...scripts: Script[]): void {
    scripts
      .filter(s => s.triggers && s.triggers.length > 0 && (!s.errors || s.errors.length === 0))
      .forEach(s => {
        this.Unregister(s)
        ;(s.triggers || [])
          .filter(t => t.eventTypes?.includes(onManual))
          .forEach(trigger => {
            const app = prop2map(trigger.uiProps).app

            if (app && this.apps.length > 0 && !this.apps.includes(app)) {
              return
            }

            try {
              const button = new Button(s, trigger)
              this.set.push(button)
              if (this.verbose) console.debug('UIHooks: registering button', s.name, { button })
            } catch (e) {
              console.warn(`UIHooks: skipping trigger of ${s.name}`, e)
            }
          })
      })

    this.set.sort(sorter)
  }

  /**
   * Removes all buttons that match a script
   */
  Unregister({ name }: Script): void {
    this.set = this.set.filter(({ script }) => name !== script)
  }

  /**
   * Searches for buttons that match the requirements
   *
   * An undefined app matches every app.
   */
  Find(resourceType: string | string[], page?: string, slot?: string, app?: string): Button[] {
    let resourceTypes: string[]

    if (!resourceType) {
      resourceTypes = []
    } else if (typeof resourceType === 'string') {
      resourceTypes = [resourceType]
    } else {
      resourceTypes = resourceType
    }

    return this.set.filter(b => {
      if (!resourceTypes.includes(b.resourceType)) {
        return false
      }

      if (app !== undefined && b.app !== app) {
        return false
      }

      return page === b.page && slot === b.slot
    })
  }

  FindByScript(script: string): Button | undefined {
    return this.set.find(b => b.script === script)
  }
}

export const UIHooksKey: InjectionKey<UIHooks> = Symbol('HumanUIHooks')

export const UIHooksPlugin: Plugin = {
  install(app: App, opts: Partial<UIHooksOptions> = {}) {
    const uiHooks = new UIHooks(opts)

    app.config.globalProperties.$UIHooks = uiHooks
    app.provide(UIHooksKey, uiHooks)
    app.provide('$UIHooks', uiHooks)
  },
}
