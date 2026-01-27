import { Module } from './module'
import { Namespace } from './namespace'
declare const fieldIndex: unique symbol
declare const propModule: unique symbol
declare const cleanValues: unique symbol
interface RawValue {
    name: string;
    value?: string;
}
interface PartialRecord extends Partial<Omit<Record, 'values' | 'createdAt' | 'updatedAt' | 'deletedAt'>> {
    values?: RawValue[];
    createdAt?: string | number | Date;
    updatedAt?: string | number | Date;
    deletedAt?: string | number | Date;
}
export interface Values {
    [name: string]: string | string[] | undefined;
}
/**
 * Combination of valid value types/structures
 */
type ValueCombo = RawValue[] | Values | Values[];
/**
 * Combination of valid types that  can be passed as Record ctor's 1st (and 2nd) parameter
 */
type RecordCtorCombo = Record | Module | PartialRecord | ValueCombo;
/**
 * Record class will be used all over the place, user scripts, etc..
 *
 * Constructor (and apply fn) is as versatile as possible to handle
 * different use-cases.
 */
export declare class Record {
  recordID: string
  moduleID: string
  namespaceID: string
  revision: number
  values: Values
  valueErrors: object
  meta: object
  createdAt?: Date
  updatedAt?: Date
  deletedAt?: Date
  ownedBy: undefined
  createdBy: undefined
  updatedBy: undefined
  deletedBy: undefined
  canUpdateRecord: boolean
  canReadRecord: boolean
  canDeleteRecord: boolean
  canUndeleteRecord: boolean
  canManageOwnerOnRecord: boolean
  canSearchRevision: boolean
  canGrant: boolean
  private [fieldIndex]
  private [propModule]?
  private [cleanValues]
  constructor(recModVal1: RecordCtorCombo, recModVal2?: RecordCtorCombo);
  clone(): Record;
  /**
     * apply (partially) updates record and it's values
     *
     * @param p
     */
  apply(p?: unknown): void;
  get cleanValues(): Values;
  get module(): Module;
  set module(m: Module);
  get namespace(): Namespace;
  /**
     * Converts internal representation of values into array of RawValue objects
     */
  serializeValues(): RawValue[];
  /**
     * Removes existing, resets default values and updates it with new ones
     */
  setValues(...i: ValueCombo[]): void;
  /**
     * Removes existing and resets default values
     */
  protected initValues(): void;
  protected defaultValues(): void;
  /**
     * Updates record's values object with provided input
     *
     * Accepted values:
     * 1. Array of RawValue objects:
     *    updateValues([{ name: ..., value: ...}, ...])
     *
     * 2. One or more Value object:
     *    updateValues({ foo: ..., bar: ... }, ...)
     */
  protected updateValues(...combo: ValueCombo[]): void;
  /**
     * Sets single value
     *
     * @param name
     * @param value
     */
  setValue(name: string, value: undefined | string | string[], index?: number): void;
  serialize(): Partial<Record>;
  /**
     * Returns resource ID
     */
  get resourceID(): string;
  /**
     * Resource type
     */
  get resourceType(): string;
  /**
     * Proxy to Record's meta to maintain BC
     */
  get labels(): object;
  get properties(): string[];
}
export {}
