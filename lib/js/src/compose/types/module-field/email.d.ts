import { ModuleField, Options } from './base'
interface EmailOptions extends Options {
  outputPlain: boolean
  multiDelimiter: string
}
export declare class ModuleFieldEmail extends ModuleField {
  readonly kind = 'Email'
  options: EmailOptions
  constructor(i?: Partial<ModuleFieldEmail>)
  applyOptions(o?: Partial<EmailOptions>): void
}
export {}
