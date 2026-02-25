import { Namespace } from '../../namespace'
import { Module } from '../../module'
import { Compose as ComposeAPI } from '../../../../api-clients'
interface FeedOptions {
  color: string
  prefilter: string
}
interface Feed {
  titleField: string
  options: FeedOptions
}
export declare function RecordFeed(
  $ComposeAPI: ComposeAPI,
  module: Module,
  namespace: Namespace,
  feed: Feed,
  options?: {},
): Promise<any[]>
export {}
