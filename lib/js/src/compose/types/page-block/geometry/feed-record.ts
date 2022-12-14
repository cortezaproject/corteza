import { Record } from '../../record'
import { Namespace } from '../../namespace'
import { Module } from '../../module'
import { Compose as ComposeAPI } from '../../../../api-clients'

interface FeedOptions {
  color: string;
  prefilter: string;
}

interface Feed {
  titleField: string;
  options: FeedOptions;
}

interface Range {
  end: Date;
  start: Date;
}

export async function RecordFeed ($ComposeAPI: ComposeAPI, module: Module, namespace: Namespace, feed: Feed): Promise<any[]> {
  // Params for record fetching
  const params = {
    namespaceID: namespace.namespaceID,
    moduleID: module.moduleID,
    query: feed.options.prefilter,
  }

  const events: Array<any> = []
  return $ComposeAPI.recordList(params).then(({ set }) => {
    return (set as Array<{ recordID: string }>)
      // cast & freeze
      .map(r => Object.freeze(new Record(module, r)))
  })
}
