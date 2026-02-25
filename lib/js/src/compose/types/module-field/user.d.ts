import { ModuleField, Options } from './base'
import { User } from '../../../system'
interface UserOptions extends Options {
  roles: Array<string>
  presetWithAuthenticated: boolean
  selectType: string
  multiDelimiter: string
  isUniqueMultiValue: boolean
}
export declare class ModuleFieldUser extends ModuleField {
  readonly kind = 'User'
  options: UserOptions
  constructor(i?: Partial<ModuleFieldUser>)
  applyOptions(o?: Partial<UserOptions>): void
  formatter({ userID, name, username, email, handle }?: Partial<User>): string
}
export {}
