import { PageBlock, PageBlockInput, Registry, BlockIssue } from '../base'
import NavigationItem, { NavigationItemInput } from './navigation-item'
import { Apply } from '../../../../cast'

const kind = 'Navigation'

interface DisplayOptions {
  appearance: string
  alignment: string
  justify: string
}

interface Options {
  display: DisplayOptions
  navigationItems: NavigationItem[]
  magnifyOption: string
}

const defaults: Readonly<Options> = Object.freeze({
  display: {
    appearance: 'pills',
    alignment: 'center',
    justify: 'justify',
  },
  navigationItems: [],
  magnifyOption: '',
})

export class PageBlockNavigation extends PageBlock {
  readonly kind = kind

  options: Options = { ...defaults }

  constructor(i?: PageBlockInput) {
    super(i)
    this.applyOptions(i?.options as Partial<Options>)
  }

  applyOptions(o?: Partial<Options>): void {
    if (!o) return

    Apply(this.options, o, String, 'magnifyOption')

    this.options.navigationItems = (o.navigationItems || []).map(f => new NavigationItem(f))

    this.options.display = { ...this.options.display, ...o.display }
  }

  static makeNavigationItem(item?: NavigationItemInput): NavigationItem {
    return new NavigationItem(item)
  }

  validate(): Array<BlockIssue> {
    const ee = super.validate()
    if (!this.options.navigationItems.length) {
      ee.push({ option: 'navigationItems', labelKey: 'block.issue.noNavigationItems' })
    }
    return ee
  }
}

Registry.set(kind, PageBlockNavigation)
