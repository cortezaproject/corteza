import { DisplayElement, DisplayElementInput, Registry } from './base'
import { Apply } from '../../../cast'

const kind = 'Text'

interface Options {
  value: string;
}

const defaults: Readonly<Options> = Object.freeze({
  value: 'Sample text...',
})

export class DisplayElementText extends DisplayElement {
  readonly kind = kind

  options: Options = { ...defaults }

  constructor (i?: DisplayElementInput) {
    super(i)
    this.applyOptions(i?.options as Partial<Options>)
  }

  applyOptions (o?: Partial<Options>): void {
    if (!o) return

    Apply(this.options, o, String, 'value')
  }
}

Registry.set(kind, DisplayElementText)
