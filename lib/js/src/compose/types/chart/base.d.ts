import { ChartConfig, Dimension, Metric, Report, TemporalDataPoint } from './util'
export type PartialChart = Partial<BaseChart>;
/**
 * BaseChart represents a structure that stores any configuration data.
 * Any display and data rendering operations should be handled by any sub classes.
 */
export declare class BaseChart {
  chartID: string
  namespaceID: string
  name: string
  handle: string
  createdAt?: Date
  updatedAt?: Date
  deletedAt?: Date
  canUpdateChart: boolean
  canDeleteChart: boolean
  canGrant: boolean
  config: ChartConfig
  constructor(def?: PartialChart);
  /**
     * The method performs post processing for each value in the given dataset.
     * It works with a simple equation written in javascript (example: n + m).
     * Available variables to use:
     * * n - current value
     * * m - previous value (undefined in case of the first element)
     * * r - entire data array.
     *
     * @param data Array of values in the given data set
     * @param m Metric for the given dataset
     */
  datasetPostProc(data: Array<number | TemporalDataPoint>, m: Metric): Array<number | TemporalDataPoint>;
  merge(c: PartialChart): void;
  /**
     * Checks reports validity.
     * Validates dimensions and metrics.
     * If invalid it throws an error.
     */
  isValid(): boolean;
  /**
     * Checks validity of dimensions.
     * If invalid it throws an error
     */
  dimCheck({ field, modifier }: Dimension): void;
  /**
     * Checks validity of metrics.
     * If invalid it throws an error
     */
  mtrCheck({ field, aggregate, type }: Metric): void;
  /**
     * Prepares params that the reporter can use for querying.
     */
  formatReporterParams({ moduleID, metrics, dimensions, filter }: Report): {
        moduleID: string | null | undefined;
        filter: string | null | undefined;
        metrics: string | undefined;
        dimensions: any;
    };
    /**
     * Fetcher reports defined in the given configuration with the help of the provided
     * reporter.
     */
  fetchReports({ reporter }: {
        reporter(p: any): Promise<any>;
    }): Promise<unknown>;
    /**
     * Processes provided report with it's results:
     * * skip missing values, if so requested,
     * * generate labels,
     * * creates dataset for the chart.
     */
  private processReporterResults
  processLabels(ll: Array<string>, d: Dimension): string[];
  makeDataset(m: Metric, d: Dimension, data: Array<number | any>, alias: string): void;
  makeOptions(data?: any): void;
  plugins(mm: Array<Metric>): void;
  baseChartType(datasets: Array<any>): void;
  /**
     * Performs chart export; used by exporter feature.
     */
  export(findModuleByID: ({ namespaceID, moduleID }: {
        namespaceID: string;
        moduleID: string;
    }) => Promise<any>): Promise<BaseChart>;
    /**
     * Performs import; used by importer feature
     */
  import(getModuleID: (moduleID: string) => string): BaseChart;
  defDimension(): Dimension;
  defMetric(): Metric;
  defReport(): Report;
  defConfig(): ChartConfig;
  /**
     * Resource type
     */
  get resourceType(): string;
  clone(): BaseChart;
}
export { chartUtil } from './util'
