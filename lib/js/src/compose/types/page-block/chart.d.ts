import { PageBlock, PageBlockInput } from './base'
import { Options as PageBlockRecordListOptions } from './record-list'
interface DrillDown {
    enabled: boolean;
    blockID: string;
    recordListOptions: Partial<PageBlockRecordListOptions>;
}
interface Options {
    chartID: string;
    refreshRate: number;
    showRefresh: boolean;
    magnifyOption: string;
    drillDown: DrillDown;
    liveFilterEnabled: boolean;
}
export declare class PageBlockChart extends PageBlock {
  readonly kind = 'Chart'
  options: Options
  constructor(i?: PageBlockInput);
  applyOptions(o?: Partial<Options>): void;
  resetDrillDown(): void;
}
export {}
