import { PluginFunction } from 'vue'
import { eventbus } from '@cortezaproject/corteza-js'

interface KV { [_: string]: string }
interface UIProp { name: string; value: string }

interface Trigger {
  resourceTypes: string[];
  eventTypes: string[];
  uiProps: UIProp[];
  constraints: object[];
  weight?: number;
}

interface Script {
  name: string;
  label: string;
  description?: string;
  errors?: string[];
  triggers: Trigger[];
}

function sorter (a: Button, b: Button): number {
  if (a.weight === b.weight) {
    return a.script.localeCompare(b.script)
  } else {
    return a.weight - b.weight
  }
}

function prop2map (uiprops: UIProp[]): KV {
  if (!uiprops) {
    return {}
  }

  return uiprops.reduce((m: KV, { name, value }) => { m[name] = value; return m }, {})
}

export class Button {
  readonly label: string
  readonly description?: string
  readonly script: string
  readonly resourceType: string
  readonly weight: number
  readonly variant?: string
  readonly page?: string
  readonly slot?: string
  readonly constraints: eventbus.ConstraintMatcher[]

  constructor (s: Script, t: Trigger) {
    const uiProps = prop2map(t.uiProps)

    if (!t.eventTypes?.includes('onManual')) {
      throw new Error('expecting onManual event type')
    }

    if (t.resourceTypes?.length !== 1) {
      throw new Error('expecting exactly one resource type on trigger')
    }

    this.label = uiProps.label ?? s.label
    this.description = s.description
    this.script = s.name
    this.weight = t.weight || 0
    this.resourceType = t.resourceTypes[0]
    this.page = uiProps.page
    this.slot = uiProps.slot
    this.variant = uiProps.variant
    this.constraints = t.constraints?.map(eventbus.ConstraintMaker) || []
  }
}

/**
 * Consumes scripts that can be triggered manually and converts it to list of buttons
 *
 * These buttons can put manually to various compose page block or
 * positioned automatically on designated pages & slots
 */
export class UIHooks {
  readonly app: string
  readonly verbose = false

  protected set: Button[] = []

  constructor (opt: string|Partial<UIHooks>) {
    if (typeof opt === 'string') {
      opt = { app: opt }
    }

    this.app = opt.app || ''
    this.verbose = !!opt?.verbose
  }

  /**
   * Takes one or more scripts and converts them to buttons
   *
   * With every script added it removes ALL
   * buttons that use the same script
   */
  Register (...scripts: Script[]): void {
    scripts
      .filter(s => s.triggers && s.triggers.length > 0 && (!s.errors || s.errors.length === 0))
      .forEach(s => {
        this.Unregister(s)

        s.triggers
          .filter(t => t.eventTypes?.includes('onManual'))
          .forEach(trigger => {
            if (prop2map(trigger.uiProps).app !== this.app) {
              // Ignore triggers that do not belong to this app.
              return
            }

            const button = new Button(s, trigger)
            this.set.push(button)
            if (this.verbose) console.debug('UIHooks: registering button', s.name, { button, trigger })
          })
      })

    // Keep buttons sorted
    this.set.sort(sorter)
  }

  /**
   * Remove all buttons that match a script
   * @param name
   * @constructor
   */
  Unregister ({ name }: Script): void {
    this.set = this.set.filter(({ script }) => name !== script)
  }

  /**
   * Searches for buttons that match the requirements
   *
   * This is used in 2 kinds of places:
   *  - currated list of buttons in compose blocks where admin can
   *    picks, reorder, name and style scripts by hand
   *  - different slots on pages where scripts are automatically placed
   *
   * @param resourceType
   * @param page
   * @param slot
   * @constructor
   */
  Find (resourceType: string|string[], page?: string, slot?: string): Button[] {
    if (!resourceType) {
      resourceType = []
    } else if (typeof resourceType === 'string') {
      resourceType = [resourceType]
    }

    resourceType = [...resourceType, 'ui:' + this.app]

    return this.set
      .filter(b => {
        if (!resourceType.includes(b.resourceType)) {
          return false
        }

        return page === b.page && slot === b.slot
      })
  }

  FindByScript (script: string): Button | undefined {
    return this.set.find(b => b.script === script)
  }
}

export default function (): PluginFunction<Partial<UIHooks>> {
  return function (Vue, opts): void {
    Vue.prototype.$UIHooks = new UIHooks(opts || {})
  }
}
