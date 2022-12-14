import { PageBlock, PageBlockInput, Registry } from '../base'
import List, { ListInput } from './list'

const kind = 'Navigation'

interface DisplayOptions {
  appearance: string;
  alignment: string;
  fillJustify: string;
}

interface Options {
  display: DisplayOptions;
  lists: List[];
}

const defaults: Readonly<Options> = Object.freeze({
  display: {
    appearance: 'tabs',
    alignment: 'left',
    fillJustify: 'justify',
  },
  lists: [],
})

export class PageBlockNavigation extends PageBlock {
  readonly kind = kind;

  options: Options = { ...defaults };

  constructor (i?: PageBlockInput) {
    super(i)
    this.applyOptions(i?.options as Partial<Options>)
  }

  applyOptions (o?: Partial<Options>): void {
    if (!o) return

    // Apply(this.options, o, String, "body");
    this.options.lists = (o.lists || []).map(f => new List(f))

    if (o.display) {
      this.options.display = o.display
    }
  }

  static makeListItem (item?: ListInput): List {
    console.log(new List(item))
    return new List(item)
  }
}

Registry.set(kind, PageBlockNavigation)
