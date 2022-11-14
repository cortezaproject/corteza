import { merge } from 'lodash'
import { IsOf } from '../../../guards'
import { Apply, CortezaID, NoID } from '../../../cast'

export const FieldNameValidator = /^[A-Za-z][0-9A-Za-z_\-.]*[A-Za-z0-9]$/

const unsortableFieldKinds = ['User', 'Record', 'File', 'Geometry']
const unsortableSysFields = ['recordID', 'ownedBy', 'createdBy', 'updatedBy', 'deletedBy']

const unfilterableFieldKinds = ['File', 'Geometry']

type fieldEncoding = null | { omit: true } | { ident: string }

export interface Capabilities {
  configurable: true;
  multi: boolean;
  writable: boolean;
  required: boolean;
  private: boolean;
}

export interface Options {
  description: {
    view: string;
    edit: string | undefined;
  };
  hint: {
    view: string;
    edit: string | undefined;
  };
}

export const defaultOptions = (): Readonly<Options> => Object.freeze({
  description: {
    view: '',
    edit: undefined,
  },
  hint: {
    view: '',
    edit: undefined,
  },
})

interface Config {
  dal: {
    encodingStrategy: fieldEncoding;
  };

  privacy: {
    sensitivityLevelID: string;
    usageDisclosure: string;
  };

  recordRevisions: {
    enabled: boolean;
  };
}

export interface Expressions {
  value?: string;

  sanitizers?: Array<string>;

  validators?: Array<Validator>;
  disableDefaultValidators?: boolean;

  formatters?: Array<string>;
  disableDefaultFormatters?: boolean;
}

interface Validator {
  validatorID: string;
  test: string;
  error: string;
}

interface DefaultValue {
  name?: string;
  value: string;
}

export class ModuleField {
  public fieldID = NoID
  public name = ''
  public kind = ''
  public label = ''

  public defaultValue: Array<DefaultValue> = []
  public maxLength = 0

  public isRequired = false
  public isMulti = false
  public isSystem = false
  public isSortable = true
  public isFilterable = true
  public options: Options = { ...defaultOptions() }
  public expressions: Expressions = {}

  public config: Partial<Config> = {
    dal: {
      encodingStrategy: null,
    },

    privacy: {
      sensitivityLevelID: NoID,
      usageDisclosure: '',
    },

    recordRevisions: {
      enabled: false,
    },
  }

  public canUpdateRecordValue = false
  public canReadRecordValue = false

  constructor (f?: Partial<ModuleField>) {
    this.apply(f)
  }

  applyOptions (o?: Partial<Options>): void {
    if (!o) return

    Apply(this.options, o, Object, 'description', 'hint')
  }

  clone (): ModuleField {
    return new ModuleField(JSON.parse(JSON.stringify(this)))
  }

  public apply (f?: Partial<ModuleField>): void {
    if (!f) return

    Apply(this, f, CortezaID, 'fieldID')
    Apply(this, f, String, 'name', 'label', 'kind')
    Apply(this, f, Number, 'maxLength')
    Apply(this, f, Boolean, 'isRequired', 'isMulti', 'isSystem')

    // Make sure field is align with it's capabilities
    if (!this.cap.multi) this.isMulti = false
    if (!this.cap.required) this.isRequired = false

    // Check if kind sortable
    if (unsortableFieldKinds.includes(this.kind)) {
      this.isSortable = false
    }

    if (unfilterableFieldKinds.includes(this.kind)) {
      this.isFilterable = false
    }

    if (f.defaultValue && Array.isArray(f.defaultValue)) {
      /**
       * Converting default value into proper format
       * so we can use it without conversion
       * false boolean values are represented only by the name, in all other cases the value is also present
       */
      this.defaultValue = f.defaultValue.filter(({ name, value }) => name !== undefined || (value !== undefined && value !== null))
    }

    if (this.isSystem) {
      this.canUpdateRecordValue = true
      this.canReadRecordValue = true
      if (unsortableSysFields.includes(this.name)) {
        this.isSortable = false
      }
    } else {
      Apply(this, f, Boolean, 'canUpdateRecordValue', 'canReadRecordValue')
    }

    if (IsOf(f, 'config')) {
      this.config = merge({}, this.config, f.config)
    }

    if (IsOf(f, 'kind')) {
      this.kind = f.kind
    }

    if (IsOf(f, 'expressions')) {
      this.expressions = f.expressions
    }
  }

  /**
   * Test field validity
   *
   * Expecting valid name
   */
  public get isValid (): boolean {
    return this.name.length > 0 && FieldNameValidator.test(this.name)
  }

  /**
   * Per module field type capabilities
   */
  public get cap (): Readonly<Capabilities> {
    return {
      configurable: true,
      multi: true,
      writable: true,
      required: true,
      private: true,
    }
  }

  /**
   * Returns resource ID
   */
  get resourceID (): string {
    return `${this.resourceType}:${this.fieldID}`
  }

  /**
   * Resource type
   */
  get resourceType (): string {
    return 'compose:module-field'
  }
}

export const Registry = new Map<string, typeof ModuleField>()
