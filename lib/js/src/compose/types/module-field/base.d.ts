export declare const FieldNameValidator: RegExp
type fieldEncoding =
  | null
  | {
      omit: true
    }
  | {
      ident: string
    }
export interface Capabilities {
  configurable: true
  multi: boolean
  writable: boolean
  required: boolean
  private: boolean
}
export interface Options {
  description: {
    view: string
    edit: string | undefined
  }
  hint: {
    view: string
    edit: string | undefined
  }
}
export declare const defaultOptions: () => Readonly<Options>
interface Config {
  dal: {
    encodingStrategy: fieldEncoding
  }
  privacy: {
    sensitivityLevelID: string
    usageDisclosure: string
  }
  recordRevisions: {
    enabled: boolean
  }
}
export interface Expressions {
  value?: string
  sanitizers?: Array<string>
  validators?: Array<Validator>
  disableDefaultValidators?: boolean
  formatters?: Array<string>
  disableDefaultFormatters?: boolean
}
interface Validator {
  validatorID: string
  test: string
  error: string
}
interface DefaultValue {
  name?: string
  value: string
}
export declare class ModuleField {
  fieldID: string
  name: string
  kind: string
  label: string
  defaultValue: Array<DefaultValue>
  maxLength: number
  isRequired: boolean
  isMulti: boolean
  isSystem: boolean
  isSortable: boolean
  isFilterable: boolean
  options: Options
  expressions: Expressions
  config: Partial<Config>
  canUpdateRecordValue: boolean
  canReadRecordValue: boolean
  constructor(f?: Partial<ModuleField>)
  applyOptions(o?: Partial<Options>): void
  clone(): ModuleField
  apply(f?: Partial<ModuleField>): void
  /**
   * Test field validity
   *
   * Expecting valid name
   */
  get isValid(): boolean
  /**
   * Per module field type capabilities
   */
  get cap(): Readonly<Capabilities>
  /**
   * Returns resource ID
   */
  get resourceID(): string
  /**
   * Resource type
   */
  get resourceType(): string
}
export declare const Registry: Map<string, typeof ModuleField>
export {}
