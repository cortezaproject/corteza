import { Validator, ValidatorFn, Validated } from '../../validator/validator'
import { Record } from '../types/record'
import { Module } from '../types/module'
export declare class RecordValidator extends Validator<Record> {
  protected rfv: {
        [field: string]: Validator<Record>;
    }
    /**
     * Construct record validator from module (or record)
     *
     * @param m
     */
  constructor(m: Module | Record);
  /**
     * Append more record validators
     *
     * @param name
     * @param vfn
     */
  push(...vfn: ValidatorFn<Record>[]): void;
  /**
     * Append more field validators
     *
     * @param name
     * @param vfn
     */
  pushToField(name: string, ...vfn: ValidatorFn<Record>[]): void;
  /**
     * Runs validators on record and all (or whitelisted) fields
     */
  run(r: Record, ...fields: string[]): Validated;
}
