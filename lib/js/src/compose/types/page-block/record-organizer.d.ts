import { PageBlock, PageBlockInput } from './base'
interface Options {
    moduleID: string;
    labelField: string;
    descriptionField: string;
    filter: string;
    positionField: string;
    groupField: string;
    group: string;
    refreshRate: number;
    showRefresh: boolean;
    magnifyOption: string;
    displayOption: string;
}
export declare class PageBlockRecordOrganizer extends PageBlock {
  readonly kind = 'RecordOrganizer'
  options: Options
  constructor(i?: PageBlockInput);
  applyOptions(o?: Partial<Options>): void;
}
export {}
