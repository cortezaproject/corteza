import { Apply } from '../../cast'
import { Param } from './param'

export interface FunctionMeta {
  short: string;
  description: string;
  visual: { [_: string]: any };
  // List of webapps where function can be used, if omitted it can be used everywhere
  webapps: Array<string>;
}

interface FunctionCtr extends Partial<Omit<Function, 'parameters' | 'results'>> {
  parameters?: Array<Partial<Param>>;
  results?: Array<Partial<Param>>;
}

export class Function {
  public ref = ''
  public kind = ''
  public meta: Partial<FunctionMeta> = {}
  public parameters: Array<Param> = []
  public results: Array<Param> = []
  public labels: { [_: string]: string } = {}

  constructor (u?: FunctionCtr) {
    this.apply(u)
  }

  apply (u?: FunctionCtr): void {
    Apply(this, u, String, 'ref', 'kind')

    if (u?.parameters) {
      this.parameters = u.parameters.map(p => new Param(p))
    }

    if (u?.results) {
      this.results = u.results.map(p => new Param(p))
    }

    if (u?.meta) {
      this.meta = { ...u.meta }
    }

    if (u?.labels) {
      this.labels = { ...u.labels }
    }
  }
}
