import { Validator, ValidatorFn, ValidatorResult, ValidatorError, Validated, IsEmpty } from '../../validator/validator'
import { Record } from '../types/record'
import { Module } from '../types/module'
import { ModuleField } from '../types/module-field'
import { IsOf } from '../../guards'
import { NoID } from '../../cast'

const emptyErr = new ValidatorError({ kind: 'empty', message: 'field:required-field' })

// Validator value types
interface FieldValidatorPayload {
  field: ModuleField;
  value: unknown | string | string[];
  oldValue: unknown | string | string[];
}

function genericFieldValidator (field: ModuleField): ValidatorFn<Record> {
  // newValue is of type unknown to satisfy ValidatorFn interface
  return function (this: Record, arg0: unknown): ValidatorResult {
    if (!IsOf<FieldValidatorPayload>(arg0, 'field', 'value', 'oldValue')) {
      throw Error('invalid field validator argument type')
    }
    const { value } = arg0

    if (field.isRequired) {
      if (value === undefined || IsEmpty(value)) {
        return emptyErr
      }
    }
  }
}

export class RecordValidator extends Validator<Record> {
  // registered field validators:
  protected rfv: { [field: string]: Validator<Record> }

  /**
   * Construct record validator from module (or record)
   *
   * @param m
   */
  constructor (m: Module|Record) {
    super()

    this.rfv = {}

    if (m instanceof Record) {
      m = m.module
    }

    m.fields.forEach(field => {
      this.rfv[field.name] = new Validator<Record>(genericFieldValidator(field))
    })
  }

  /**
   * Append more record validators
   *
   * @param name
   * @param vfn
   */
  public push (...vfn: ValidatorFn<Record>[]): void {
    this.registered.push(...vfn)
  }

  /**
   * Append more field validators
   *
   * @param name
   * @param vfn
   */
  public pushToField (name: string, ...vfn: ValidatorFn<Record>[]): void {
    if (!this.rfv[name]) {
      throw new Error('can not push validators to unknown field')
    }

    this.rfv[name].push(...vfn)
  }

  /**
   * Runs validators on record and all (or whitelisted) fields
   */
  public run (r: Record, ...fields: string[]): Validated {
    const out = new Validated()

    if (fields.length === 0) {
      // Fields are not explicitly provided,
      // we can run record-wide validators:
      const result = super.run(r)
      out.push(result.get())

      // get list of fields from registered field validators
      fields = Object.getOwnPropertyNames(this.rfv)
    }

    for (const f of fields) {
      const field = r.module.findField(f)
      if (!field) {
        continue
      }

      const payload: FieldValidatorPayload = {
        value: r.values[f],
        oldValue: r.cleanValues[f],
        field,
      }

      const { recordID = NoID } = r || {}
      const results = this.rfv[f].run(r, payload)

      results.applyMeta({ field: f, id: recordID === NoID ? 'parent:0' : recordID })
      out.push(results.get())
    }

    return out
  }
}
