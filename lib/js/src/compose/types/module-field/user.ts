// @todo option to allow multiple entries
// @todo option to allow duplicates
import { ModuleField, Registry, Options, defaultOptions } from './base'
import { Apply } from '../../../cast'
import { User } from '../../../system'

const kind = 'User'

interface UserOptions extends Options {
  roles: Array<string>;
  presetWithAuthenticated: boolean;
  selectType: string;
  multiDelimiter: string;
}

const defaults = (): Readonly<UserOptions> => Object.freeze({
  ...defaultOptions(),
  roles: [],
  presetWithAuthenticated: false,
  selectType: 'default',
  multiDelimiter: '\n',
})

export class ModuleFieldUser extends ModuleField {
  readonly kind = kind

  options: UserOptions = { ...defaults() }

  constructor (i?: Partial<ModuleFieldUser>) {
    super(i)
    this.applyOptions(i?.options)
  }

  applyOptions (o?: Partial<UserOptions>): void {
    if (!o) return
    super.applyOptions(o)

    Apply(this.options, o, Boolean, 'presetWithAuthenticated')
    Apply(this.options, o, String, 'selectType', 'multiDelimiter')
    Apply(this.options, o, (o) => {
      if (!o) {
        return []
      }
      if (!Array.isArray(o)) {
        return [o]
      }
      return o
    }, 'roles')
  }

  formatter ({ userID, name, username, email, handle }: Partial<User> = {}): string {
    return name || username || email || handle || userID || ''
  }
}

Registry.set(kind, ModuleFieldUser)
