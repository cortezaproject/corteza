import { Record } from '../../record'
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

export async function RecordFeed(
  $ComposeAPI: ComposeAPI,
  module: Module,
  namespace: Namespace,
  feed: Feed,
  options = {},
): Promise<any[]> {
  // Params for record fetching
  const params = {
    namespaceID: namespace.namespaceID,
    moduleID: module.moduleID,
    query: feed.options.prefilter,
  }

  return $ComposeAPI.recordList(params, options).then(({ set }) => {
    return (
      (set as Array<{ recordID: string }>)
        // cast & freeze
        .map(r => Object.freeze(new Record(module, r)))
    )
  })
}
