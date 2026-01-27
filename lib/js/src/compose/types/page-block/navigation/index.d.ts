import { PageBlock, PageBlockInput } from '../base'
import NavigationItem, { NavigationItemInput } from './navigation-item'
interface DisplayOptions {
    appearance: string;
    alignment: string;
    justify: string;
}
interface Options {
    display: DisplayOptions;
    navigationItems: NavigationItem[];
    magnifyOption: string;
}
export declare class PageBlockNavigation extends PageBlock {
  readonly kind = 'Navigation'
  options: Options
  constructor(i?: PageBlockInput);
  applyOptions(o?: Partial<Options>): void;
  static makeNavigationItem(item?: NavigationItemInput): NavigationItem;
}
export {}
