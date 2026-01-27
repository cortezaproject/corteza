import { ModuleField, Options } from './base'
interface SelectOptionStyle {
    textColor?: string;
    backgroundColor?: string;
}
interface SelectOption {
    value: string;
    text: string;
    style: SelectOptionStyle;
}
interface SelectOptions extends Options {
    options: Array<SelectOption>;
    selectType: string;
    displayType: 'text' | 'badge';
    multiDelimiter: string;
    isUniqueMultiValue: boolean;
}
export declare class ModuleFieldSelect extends ModuleField {
  readonly kind = 'Select'
  options: SelectOptions
  constructor(i?: Partial<ModuleFieldSelect>);
  applyOptions(o?: Partial<SelectOptions>): void;
  createSelectOption({ value, text, style }?: Partial<SelectOption>): SelectOption;
}
export {}
