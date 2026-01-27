export type PartialDisplayElement = Partial<DisplayElement>;
export type DisplayElementInput = DisplayElement | PartialDisplayElement;
export declare class DisplayElement {
  elementID: string
  name: string
  description: string
  options: {}
  meta: {
        size: undefined;
    }
  kind: string
  constructor(de?: PartialDisplayElement);
  apply(de?: DisplayElement | PartialDisplayElement): void;
}
export declare const Registry: Map<string, typeof DisplayElement>
