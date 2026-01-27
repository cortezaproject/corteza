interface PageBlockStyleVariants {
    [_: string]: string;
}
interface PageBlockStyleWrap {
    kind: string;
}
interface PageBlockStyleBorder {
    enabled?: boolean;
}
interface PageBlockStyle {
    variants: PageBlockStyleVariants;
    wrap?: PageBlockStyleWrap;
    border?: PageBlockStyleBorder;
}
interface Visibility {
    expression: string;
    roles: string[];
}
interface PageBlockMeta {
    hidden?: boolean;
    tempID?: string;
    customID?: string;
    customCSSClass?: string;
    visibility: Visibility;
}
export type PageBlockInput = PageBlock | Partial<PageBlock>;
export declare class PageBlock {
  blockID: string
  kind: string
  title: string
  description: string
  xywh: number[]
  options: {}
  meta: PageBlockMeta
  style: PageBlockStyle
  constructor(i?: PageBlockInput);
  apply(i?: PageBlockInput): void;
  validate(): Array<string>;
  setTempID(): void;
  clone(): this;
}
export declare const Registry: Map<string, typeof PageBlock>
export {}
