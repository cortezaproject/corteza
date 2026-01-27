import { ModuleField, Options } from './base'
interface UrlOptions extends Options {
    trimFragment: boolean;
    trimQuery: boolean;
    trimPath: boolean;
    onlySecure: boolean;
    outputPlain: boolean;
    multiDelimiter: string;
}
export declare class ModuleFieldUrl extends ModuleField {
  readonly kind = 'Url'
  options: UrlOptions
  constructor(i?: Partial<ModuleFieldUrl>);
  applyOptions(o?: Partial<UrlOptions>): void;
}
export {}
