import { PageBlock, Registry } from './base'
export { PageBlockAutomation } from './automation'
export { PageBlockChart } from './chart'
export { PageBlockContent } from './content'
export { PageBlockFile } from './file'
export { PageBlockIFrame } from './iframe'
export { PageBlockRecord } from './record'
export { PageBlockRecordList } from './record-list'
export { PageBlockRecordRevisions } from './record-revisions'
export { PageBlockRecordOrganizer } from './record-organizer'
export { PageBlockSocialFeed } from './social-feed'
export { PageBlockCalendar } from './calendar'
export { PageBlockMetric } from './metric'
export { PageBlockComment } from './comment'
export { PageBlockReport } from './report'
export { PageBlockProgress } from './progress'
export { PageBlockNavigation } from './navigation'
export { PageBlockTab } from './tabs'
export { PageBlockGeometry } from './geometry'

export function PageBlockMaker<T extends PageBlock> (i: { kind: string }): T {
  const PageBlockTemp = Registry.get(i.kind)
  if (PageBlockTemp === undefined) {
    throw new Error(`unknown block kind '${i.kind}'`)
  }

  if (i instanceof PageBlock) {
    // Get rid of the references
    i = JSON.parse(JSON.stringify(i))
  }

  return new PageBlockTemp(i) as T
}

export {
  Registry as PageBlockRegistry,
  PageBlock,
}
