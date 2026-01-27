export interface Typed {
    '@type': string;
    '@value': any;
}
export type Vars = {
    [_: string]: Typed;
};
export declare function IsTyped(a: unknown): a is Typed;
/**
 *
 * @param any
 * @constructor
 */
export declare function Encode(input: {
    [_: string]: any;
}): Vars;
