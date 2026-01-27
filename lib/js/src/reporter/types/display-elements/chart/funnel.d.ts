import { ChartOptions } from './base'
import { FrameDefinition } from '../../frame'
export declare class FunnelChartOptions extends ChartOptions {
  labelColumn: string
  dataColumns: Array<{
        name: string;
        label?: string;
    }>
  constructor(o?: FunnelChartOptions | Partial<FunnelChartOptions>);
  getChartConfiguration(dataframes: Array<FrameDefinition>, meta: any): {
        animation: boolean;
        title: {
            text: string;
            left: string;
            textStyle: {
                fontFamily: any;
                color: any;
                fontSize: number;
            };
        };
        textStyle: {
            fontFamily: any;
        };
        tooltip: {
            show: boolean;
            trigger: string;
            formatter: (params: any) => string;
            appendToBody: boolean;
        };
        legend: {
            show: boolean;
            type: string;
            top: string | number | undefined;
            right: string | undefined;
            bottom: string | undefined;
            left: string;
            orient: string;
            textStyle: {
                color: any;
            };
            pageTextStyle: {
                color: any;
            };
            pageIconColor: any;
            pageIconInactiveColor: any;
        };
        series: {
            type: string;
            name: string;
            sort: string;
            width: string;
            label: {
                show: boolean;
                position: string;
                align: string;
                verticalAlign: string;
                formatter: string;
            };
            data: {
                name: string;
                value: any;
                itemStyle: {
                    color: string;
                };
            }[];
            top: string | number | undefined;
            right: string | undefined;
            bottom: string | undefined;
            left: string | undefined;
        }[];
    };
  getColIndex(dataframe: FrameDefinition, col: string): number;
  getLabels(localDataframe: FrameDefinition): string[];
  getDatasets(localDataframe: FrameDefinition, dataframes: Array<FrameDefinition>): any;
}
