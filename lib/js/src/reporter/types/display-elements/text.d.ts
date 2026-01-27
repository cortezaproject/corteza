import { DisplayElement, DisplayElementInput } from './base'
interface Options {
    value: string;
}
export declare class DisplayElementText extends DisplayElement {
  readonly kind = 'Text'
  options: Options
  constructor(i?: DisplayElementInput);
  applyOptions(o?: Partial<Options>): void;
}
export {}
