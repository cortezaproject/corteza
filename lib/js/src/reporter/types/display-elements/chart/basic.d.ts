import { ChartOptions } from './base'
import { FrameDefinition } from '../../frame'
export declare class BasicChartOptions extends ChartOptions {
  labelColumn: string
  dataColumns: Array<{
        name: string;
        label?: string;
        stack?: string;
    }>
  constructor(o?: BasicChartOptions | Partial<BasicChartOptions>);
  getChartConfiguration(dataframes: Array<FrameDefinition>, meta: any): any;
  getColIndex(dataframe: FrameDefinition, col: string): number;
  getData(localDataframe: FrameDefinition, dataframes: Array<FrameDefinition>): {
        datasets: any[];
        labels: string[];
    };
}
