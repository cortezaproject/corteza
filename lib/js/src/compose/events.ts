import { Event, onManual } from '../eventbus/shared'
import { Module } from './types/module'
import { Page } from './types/page'
import { Record, Values } from './types/record'
import { Namespace } from './types/namespace'
import { ConstraintMatcher } from '../eventbus/constraints'
import { IsOf } from '../guards'

interface TriggerEndpoints {
  automationTriggerScript (params: { script: string }): Promise<object>;
  namespaceTriggerScript (params: { namespaceID: string; script: string }): Promise<object>;
  moduleTriggerScript (params: { namespaceID: string; moduleID: string; script: string }): Promise<object>;
  recordTriggerScript (params: { namespaceID: string; moduleID: string; recordID: string; values: Values; script: string }): Promise<object>;
}

function namespaceMatcher (r: Namespace, c: ConstraintMatcher, def: boolean): boolean {
  if (!r) {
    throw new Error('can not run namespace matcher on undefined/null namespace')
  }

  // keep in sync with server/compose/service/event/namespace.go
  switch (c.Name()) {
    case 'namespace':
    case 'namespace.slug':
      return c.Match(r.slug)
    case 'namespace.name':
      return c.Match(r.name)
  }

  return def
}

function moduleMatcher (r: Module, c: ConstraintMatcher, def: boolean): boolean {
  if (!r) {
    throw new Error('can not run module matcher on undefined/null module')
  }

  // keep in sync with server/compose/service/event/module.go
  switch (c.Name()) {
    case 'module':
    case 'module.handle':
      return c.Match(r.handle)
    case 'module.name':
      return c.Match(r.name)
  }

  return def
}

/**
 * Creates event for compose resource with ready-to-go-defaults
 */
export function ComposeEvent (event?: Partial<Event>): Event {
  return {
    eventType: onManual,
    resourceType: 'compose',
    match: (): boolean => true,
    ...event,
  }
}

/**
 * Creates namespace event with ready-to-go-defaults
 */
export function NamespaceEvent (res: Namespace, event?: Partial<Event>): Event {
  return {
    eventType: onManual,
    resourceType: res.resourceType,
    match: (c): boolean => namespaceMatcher(res, c, false),
    ...event,

    // Override the arguments at the end
    args: { namespace: res, ...event?.args },
  }
}

/**
 * Creates module event with ready-to-go-defaults
 */
export function ModuleEvent (res: Module, event?: Partial<Event>): Event {
  return {
    eventType: onManual,
    resourceType: res.resourceType,
    match: (c): boolean => namespaceMatcher(res.namespace, c, moduleMatcher(res, c, false)),
    ...event,

    // Override the arguments at the end
    args: { module: res, namespace: res.namespace, ...event?.args },
  }
}

/**
 * Creates record event with ready-to-go-defaults
 */
export function RecordEvent (res: Record, event?: Partial<Event>): Event {
  return {
    eventType: onManual,
    resourceType: res.resourceType,
    match: (c): boolean => namespaceMatcher(res.namespace, c, moduleMatcher(res.module, c, false)),
    ...event,

    // Override the arguments at the end
    args: { record: res, module: res.module, namespace: res.namespace, ...event?.args },
  }
}

/**
 * Creates record event with ready-to-go-defaults
 */
export function PageEvent (res: Page, event?: Partial<Event>): Event {
  return {
    eventType: onManual,
    resourceType: 'compose:page',
    match: (): boolean => true,
    ...event,

    // Override the arguments at the end
    args: { page: res, ...event?.args },
  }
}

/**
 * Returns handler that routes onManual events for server script to the compose API
 *
 * See makeAutomationScriptsRegistrator
 *
 * @param api
 * @return function
 */
export function TriggerComposeServerScriptOnManual (api: TriggerEndpoints) {
  return (ev: Event, script: string): Promise<unknown> => {
    const params = { script, args: ev.args }

    if (ev.resourceType === 'compose') {
      return api
        .automationTriggerScript({ ...params })
    }

    if (!ev.args) {
      throw new Error('expecting args prop in event')
    }

    if (ev.resourceType === 'compose:namespace') {
      if (!IsOf<Namespace>(ev.args.namespace, 'namespaceID')) {
        throw new Error('expecting args.namespace in event arguments')
      }

      const { namespaceID } = ev.args.namespace as Namespace
      return api
        .namespaceTriggerScript({ namespaceID, ...params })
    }

    if (ev.resourceType === 'compose:module') {
      if (!IsOf<Module>(ev.args.module, 'namespaceID', 'moduleID')) {
        throw new Error('expecting args.module in event arguments')
      }

      const { namespaceID, moduleID } = ev.args.module as Module

      return api
        .moduleTriggerScript({ namespaceID, moduleID, ...params })
    }

    if (ev.resourceType === 'compose:record') {
      if (!IsOf<Record>(ev.args.record, 'namespaceID', 'moduleID', 'recordID')) {
        throw new Error('expecting args.record in event arguments')
      }

      const record = ev.args.record as Record
      const { namespaceID, moduleID, recordID, values } = record

      return api
        .recordTriggerScript({ namespaceID, moduleID, recordID, values, ...params })
        .then(rval => record.apply(rval))
    }

    throw Error(`cannot trigger server script: unknown resource type '${ev.resourceType}'`)
  }
}
