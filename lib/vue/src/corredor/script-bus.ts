import { eventbus } from '@planetcrust/human-js'
import type { App, InjectionKey, Plugin } from 'vue'

const onManual = 'onManual'

/**
 * Event types a record page dispatches around a save, a delete and a restore.
 */
const recordPageEventTypes = [
  'beforeFormSubmit',
  'onFormSubmitError',
  'afterFormSubmit',
  'beforeDelete',
  'afterDelete',
  'beforeUndelete',
  'afterUndelete',
]

/**
 * Resource & event type pairs the script bus knows about.
 *
 * `ui:` resources are dispatched by the webapp itself; the rest mirror the
 * resources Corredor scripts can bind a trigger to.
 */
export const wellKnownScriptPairs: eventbus.WellKnownPairs = {
  compose: [onManual],
  'compose:namespace': [onManual],
  'compose:module': [onManual],
  'compose:record': [onManual],
  'compose:page': [onManual],
  'ui:compose': [onManual],
  'ui:compose:record-page': [onManual, ...recordPageEventTypes],
  'ui:compose:admin-record-page': [onManual, ...recordPageEventTypes],
  system: [onManual],
  'system:user': [onManual],
  'system:role': [onManual],
}

export const ScriptBusKey: InjectionKey<eventbus.EventBus> = Symbol('HumanScriptBus')

/**
 * Bus Corredor scripts are registered on and webapp events dispatched to.
 */
export const ScriptBusPlugin: Plugin = {
  install(app: App) {
    const scriptBus = new eventbus.EventBus({
      pairs: wellKnownScriptPairs,
      strict: false,
      verbose: false,
    })

    app.config.globalProperties.$ScriptBus = scriptBus
    app.provide(ScriptBusKey, scriptBus)
    app.provide('$ScriptBus', scriptBus)
  },
}
