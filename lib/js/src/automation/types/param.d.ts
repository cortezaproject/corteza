export interface ParamMeta {
    label: string;
    description: string;
    visual: {
        [_: string]: any;
    };
}
export declare class Param {
  name: string
  types: Array<string>
  required: boolean
  isArray: boolean
  meta: Partial<ParamMeta>
  constructor(u?: Partial<Param>);
  apply(u?: Partial<Param>): void;
}
