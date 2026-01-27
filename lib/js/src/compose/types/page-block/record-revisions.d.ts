import { PageBlock, PageBlockInput } from './base'
import { Compose as ComposeAPI } from '../../../api-clients'
import { Record } from '../../'
import { Revision } from '../revision'
interface Options {
    preload: boolean;
    displayedFields: string[];
    expRefFields: string[];
    refreshRate: number;
    showRefresh: boolean;
    magnifyOption: string;
}
export declare class PageBlockRecordRevisions extends PageBlock {
  readonly kind = 'RecordRevisions'
  options: Options
  constructor(i?: PageBlockInput);
  applyOptions(o?: Partial<Options>): void;
  /**
     * fetch is a utility method on record revision page block
     * that fetches revisions for a record and converts them to RevisionPayload class
     *
     * this function also strips out all fields that should not be dispalyed
     * (as per displayedFields option)
     *
     * @param api Compose API to be used
     * @param record Record to fetch revisions for
     */
  fetch(api: ComposeAPI, record: Record): Promise<Array<Revision>>;
}
export {}
