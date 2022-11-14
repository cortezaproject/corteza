import { Apply } from '../../cast'

export interface ParamMeta {
  label: string;
  description: string;
  visual: { [_: string]: any };
}

export class Param {
  public name = ''
  public types: Array<string> = []
  public required = false
  public isArray = false
  public meta: Partial<ParamMeta> = {}

  constructor (u?: Partial<Param>) {
    this.apply(u)
  }

  apply (u?: Partial<Param>): void {
    Apply(this, u, String, 'name')
    Apply(this, u, Boolean, 'required', 'isArray')

    if (u?.types) {
      this.types = u.types
    }

    if (u?.meta) {
      this.meta = { ...u.meta }
    }
  }
}
