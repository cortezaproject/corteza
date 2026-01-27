import { Param } from './param'
export interface FunctionMeta {
    short: string;
    description: string;
    visual: {
        [_: string]: any;
    };
    webapps: Array<string>;
}
interface FunctionCtr extends Partial<Omit<Function, 'parameters' | 'results'>> {
    parameters?: Array<Partial<Param>>;
    results?: Array<Partial<Param>>;
}
export declare class Function {
  ref: string
  kind: string
  meta: Partial<FunctionMeta>
  parameters: Array<Param>
  results: Array<Param>
  labels: {
        [_: string]: string;
    }
  constructor(u?: FunctionCtr);
  apply(u?: FunctionCtr): void;
}
export {}
