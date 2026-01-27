import { ModuleField, Options } from './base'
interface Threshold {
    value: number;
    variant: string;
}
interface NumberOptions extends Options {
    presetFormat: string;
    format: string;
    prefix: string;
    suffix: string;
    precision: number;
    multiDelimiter: string;
    display: string;
    min: number;
    max: number;
    showValue: boolean;
    showRelative: boolean;
    showProgress: boolean;
    animated: boolean;
    variant: string;
    thresholds: Threshold[];
}
export declare class ModuleFieldNumber extends ModuleField {
  readonly kind = 'Number'
  options: NumberOptions
  constructor(i?: Partial<ModuleFieldNumber>);
  applyOptions(o?: Partial<NumberOptions>): void;
  formatValue(value: string, format: string): string;
}
export {}
