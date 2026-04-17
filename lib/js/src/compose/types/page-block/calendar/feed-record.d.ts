import { Namespace } from '../../namespace'
import { Module } from '../../module'
import { Compose as ComposeAPI } from '../../../../api-clients'
import { Event } from './shared'
interface FeedOptions {
  color: string
  prefilter: string
}
interface Feed {
  startField: string
  endField: string
  titleField: string
  options: FeedOptions
  allDay: boolean
}
interface Range {
  end: Date
  start: Date
}
/**
 * Loads & converts module resource into FC events
 * @param {ComposeAPI} $ComposeAPI ComposeAPI provider
 * @param {Module} module Current module
 * @param {Namespace} namespace Current namespace
 * @param {Feed} feed Current feed
 * @param {Object} range Current date range
 * @returns {Promise<Array>} Resolves to a set of FC events to display
 */
export declare function RecordFeed(
  $ComposeAPI: ComposeAPI,
  module: Module,
  namespace: Namespace,
  feed: Feed,
  range: Range,
  options?: Record<string, unknown>,
): Promise<Event[]>
export {}
