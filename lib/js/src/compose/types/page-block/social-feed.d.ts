import { PageBlock, PageBlockInput } from './base'
interface Options {
    moduleID: string;
    fields: unknown[];
    profileSourceField: string;
    profileUrl: string;
    showRefresh: boolean;
    refreshRate: number;
    magnifyOption: string;
}
export declare class PageBlockSocialFeed extends PageBlock {
  readonly kind = 'SocialFeed'
  options: Options
  constructor(i?: PageBlockInput);
  applyOptions(o?: Partial<Options>): void;
}
export {}
