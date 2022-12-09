import { PageBlock, PageBlockInput, Registry } from '../base'
import NavigationItem, { NavigationItemInput } from './navigation-item'

const kind = 'Navigation'

interface DisplayOptions {
  appearance: string;
  alignment: string;
  fillJustify: string;
}

interface Options {
  display: DisplayOptions;
  navigationItems: NavigationItem[];
}

const defaults: Readonly<Options> = Object.freeze({
  display: {
    appearance: 'tabs',
    alignment: 'center',
    fillJustify: 'justify',
  },
  navigationItems: [],
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

    this.options.navigationItems = (o.navigationItems || []).map(f => new NavigationItem(f))

    if (o.display) {
      this.options.display = o.display
    }
  }

  static makeNavigationItem (item?: NavigationItemInput): NavigationItem {
    return new NavigationItem(item)
  }
}

Registry.set(kind, PageBlockNavigation)
