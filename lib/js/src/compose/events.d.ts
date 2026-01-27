import { Event } from '../eventbus/shared'
import { Module } from './types/module'
import { Namespace } from './types/namespace'
import { Page } from './types/page'
import { Record, Values } from './types/record'
interface TriggerEndpoints {
  automationTriggerScript(params: { script: string }): Promise<object>
  namespaceTriggerScript(params: { namespaceID: string; script: string }): Promise<object>
  moduleTriggerScript(params: {
    namespaceID: string
    moduleID: string
    script: string
  }): Promise<object>
  recordTriggerScript(params: {
    namespaceID: string
    moduleID: string
    recordID: string
    values: Values
    script: string
  }): Promise<object>
}
/**
 * Creates event for compose resource with ready-to-go-defaults
 */
export declare function ComposeEvent(event?: Partial<Event>): Event
/**
 * Creates namespace event with ready-to-go-defaults
 */
export declare function NamespaceEvent(res: Namespace, event?: Partial<Event>): Event
/**
 * Creates module event with ready-to-go-defaults
 */
export declare function ModuleEvent(res: Module, event?: Partial<Event>): Event
/**
 * Creates record event with ready-to-go-defaults
 */
export declare function RecordEvent(res: Record, event?: Partial<Event>): Event
/**
 * Creates record event with ready-to-go-defaults
 */
export declare function PageEvent(res: Page, event?: Partial<Event>): Event
/**
 * Returns handler that routes onManual events for server script to the compose API
 *
 * See makeAutomationScriptsRegistrator
 *
 * @param api
 * @return function
 */
export declare function TriggerComposeServerScriptOnManual(
  api: TriggerEndpoints,
): (ev: Event, script: string) => Promise<unknown>
export {}
