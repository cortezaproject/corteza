import { PageBlock, PageBlockInput } from './base'
import { Button } from './types'
interface Options {
    buttons: Array<Button>;
    sealed: boolean;
    magnifyOption: string;
}
export declare class PageBlockAutomation extends PageBlock {
  readonly kind = 'Automation'
  options: Options
  constructor(i?: PageBlockInput);
  applyOptions(o?: Partial<Options>): void;
  validate(): Array<string>;
}
export {}
