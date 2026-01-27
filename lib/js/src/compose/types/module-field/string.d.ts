import { ModuleField, Options } from './base'
interface StringOptions extends Options {
    multiLine: boolean;
    useRichTextEditor: boolean;
    multiDelimiter: string;
}
export declare class ModuleFieldString extends ModuleField {
  readonly kind = 'String'
  options: StringOptions
  constructor(i?: Partial<ModuleFieldString>);
  applyOptions(o?: Partial<StringOptions>): void;
}
export {}
