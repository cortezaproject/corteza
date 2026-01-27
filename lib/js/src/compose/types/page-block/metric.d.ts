import { PageBlock, PageBlockInput } from './base'
import { Options as PageBlockRecordListOptions } from './record-list'
type Reporter = (p: ReporterParams) => Promise<any>;
interface DrillDown {
    enabled: boolean;
    blockID: string;
    recordListOptions: Partial<PageBlockRecordListOptions>;
}
interface ReporterParams {
    moduleID: string;
    filter?: string;
    metrics?: string;
    dimensions: string;
}
interface Style {
    color: string;
    backgroundColor: string;
    fontSize?: string;
}
interface Metric {
    label: string;
    moduleID: string;
    dimensionField: string;
    dateFormat?: string;
    filter?: string;
    bucketSize?: string;
    metricField: string;
    operation: string;
    numberFormat?: string;
    prefix?: string;
    suffix?: string;
    transformFx?: string;
    valueStyle?: Style;
    drillDown: DrillDown;
}
interface Options {
    metrics: Array<Metric>;
    refreshRate: number;
    showRefresh: boolean;
    magnifyOption: string;
}
export declare class PageBlockMetric extends PageBlock {
  readonly kind = 'Metric'
  options: Options
  constructor(i?: PageBlockInput);
  applyOptions(o?: Partial<Options>): void;
  /**
     * Helper function to fetch and parse reporter's reports.
     */
  fetch({ m }: {
        m: Metric;
    }, reporter: Reporter): Promise<object>;
    /**
     * Helper to construct reporter's params
     */
  private formatParams
  makeMetric(): Readonly<Metric>;
}
export {}
