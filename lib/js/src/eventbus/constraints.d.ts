interface Constraint {
    name?: string;
    op?: string;
    value: string[];
}
export interface ConstraintMatcher {
    Name(): string | undefined;
    Values(): string[];
    Match(value: string): boolean;
}
export declare class Equal {
  readonly name?: string
  readonly values: string[]
  protected not: boolean
  constructor(name: string | undefined, vv: string[], not?: boolean);
  Name(): string | undefined;
  Values(): string[];
  Match(value: string): boolean;
}
/**
 * Handle glob-like pattern matching
 *
 * See: https://github.com/isaacs/minimatch
 */
export declare class Like extends Equal {
  constructor(name: string | undefined, vv: string[], not?: boolean);
  Match(value: string): boolean;
}
/**
 * Regex matcher
 */
export declare class Match extends Equal {
  protected re: RegExp[]
  constructor(name: string | undefined, vv: string[], not?: boolean);
  Match(value: string): boolean;
}
export declare function ConstraintMaker(c: Constraint | unknown): ConstraintMatcher;
export {}
