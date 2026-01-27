import { Capabilities, ModuleField, Options } from './base'
interface BoolOptions extends Options {
    trueLabel: string;
    falseLabel: string;
    switch: boolean;
}
export declare class ModuleFieldBool extends ModuleField {
  readonly kind = 'Bool'
  options: BoolOptions
  constructor(i?: Partial<ModuleFieldBool>);
  applyOptions(o?: Partial<BoolOptions>): void;
  /**
     * Per module field type capabilities
     */
  get cap(): Readonly<Capabilities>;
}
export {}
