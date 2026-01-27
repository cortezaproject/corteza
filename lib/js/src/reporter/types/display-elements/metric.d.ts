import { DisplayElement, DisplayElementInput } from './base'
import { FrameDefinition, DefinitionOptions } from '../frame'
interface Options {
    source?: string;
    datasources: Array<FrameDefinition>;
    valueColumn: string;
    format: string;
    prefix: string;
    suffix: string;
    color: string;
    backgroundColor: string;
}
export declare class DisplayElementMetric extends DisplayElement {
  readonly kind = 'Metric'
  options: Options
  constructor(i?: DisplayElementInput);
  applyOptions(o?: Partial<Options>): void;
  reportDefinitions(definition?: DefinitionOptions): {
        dataframes: Array<FrameDefinition>;
    };
}
export {}
