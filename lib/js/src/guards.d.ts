export declare const IsOf: <T>(v: unknown, ...props: (keyof T)[]) => v is T
export declare const AreStrings: (a: unknown | unknown[]) => a is string[]
export declare const AreBooleans: (a: unknown | unknown[]) => a is boolean[]
export declare const AreNumbers: (a: unknown | unknown[]) => a is number[]
export declare const AreObjects: (a: unknown | unknown[]) => a is object[]
export declare function AreObjectsOf<T>(a: unknown | unknown[], ...props: (keyof T)[]): a is T[];
