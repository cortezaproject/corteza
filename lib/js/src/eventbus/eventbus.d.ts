/**
 * EventBus handles all Corteza events on browser (!! not on corredor server !!)
 *
 * Flow #1
 *  1. Corredor prepares a bundle that is loaded on a client
 *  2. Bundle provides a "callback" function that accepts EventBus object +
 *     all context information and configuration that is needed for
 *     handler registration
 * 3a. When a "Corteza event" is dispatched (Dispatch())
 *     event-bus searches for handler and passes on event information
 * 3b. When manual event is executed (Exec())
 *     event-bus searches for handler and passes on event information
 *
 * Flow #2
 *   1. When web application is initialized, it should register all
 *      explicit server scripts
 *   2. These server scripts are wrapped with a handlerFn that forwards call
 *      to the API (there, request is passed to the Corredor where it's executed)
 *
 */
import { Event, HandlerFn, Trigger } from './shared'
export interface WellKnownPairs {
    [resource: string]: Array<string>;
}
export interface Options {
    pairs: WellKnownPairs;
    strict: boolean;
    verbose: boolean;
}
/**
 * EventBus for event dispatching and handling
 *
 * Since we have much shorter execution path here than we have in case of server scripts,
 * we can afford some optimisation (in comparison to backend's pkg/eventbus)
 */
export declare class EventBus {
  /**
     * List of wellknown resource & event type pairs
     *
     * If set, eventbus will throw error if unresognized pair is registered or dispatched
     */
  readonly pairs?: WellKnownPairs
  readonly strict: boolean
  readonly verbose: boolean
  private handlers
  constructor(opt?: Partial<Options>);
  /**
     * Dispatches event and sequentially calls all handlers.
     *
     * Handling handler results works a bit different then on backend.
     * Scripts executed with handlers have DIRECT access to values passed (by reference)
     * as arguments via event so there's no need to do an explicit return
     *
     * @param {Event} ev Event to dispatch
     */
  Dispatch(ev: Event, script?: string): Promise<null>;
  /**
     * Filters and sorts all handlers by event & constraints
     */
  private find
  /**
     * Registers Event handler
     *
     * @param handler Handler function
     * @param trigger Trigger definition
     */
  Register(handler: HandlerFn, trigger: Trigger): EventBus;
  /**
     * Unregisters all handlers
     */
  UnregisterAll(): EventBus;
  protected checkPairs(resourceTypes: string[], eventTypes: string[]): void;
}
