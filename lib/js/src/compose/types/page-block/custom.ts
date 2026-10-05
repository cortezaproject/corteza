import { PageBlock, PageBlockInput, Registry, BlockIssue } from './base'
import { Apply } from '../../../cast'

const kind = 'Custom'

// A custom HTML block: either a custom application's page, by applicationID,
// or a page of its own held in `source`, with the modules of the page's
// namespace it may read and change beside it.
interface Options {
  applicationID: string
  source: string
  modules: Array<string>
  // Each module handle beside its ID, as the server stored them.
  moduleIDs: Record<string, string>
  writes: Array<string>
  origins: Array<string>
  automations: Array<string>
  chatbots: Array<string>
  params: Record<string, unknown>
  magnifyOption: string
}

const defaults: Readonly<Options> = Object.freeze({
  applicationID: '',
  source: '',
  modules: [],
  moduleIDs: {},
  writes: [],
  origins: [],
  automations: [],
  chatbots: [],
  params: {},
  magnifyOption: '',
})

function strings(v: unknown): Array<string> {
  return Array.isArray(v) ? v.map(String).filter(s => !!s) : []
}

export class PageBlockCustom extends PageBlock {
  readonly kind = kind

  options: Options = {
    ...defaults,
    modules: [],
    moduleIDs: {},
    writes: [],
    origins: [],
    automations: [],
    chatbots: [],
    params: {},
  }

  constructor(i?: PageBlockInput) {
    super(i)
    this.applyOptions(i?.options as Partial<Options>)
  }

  applyOptions(o?: Partial<Options>): void {
    if (!o) return

    Apply(this.options, o, String, 'applicationID', 'source', 'magnifyOption')
    if (o.modules !== undefined) this.options.modules = strings(o.modules)
    if (o.moduleIDs !== undefined) {
      this.options.moduleIDs =
        o.moduleIDs && typeof o.moduleIDs === 'object' ? { ...o.moduleIDs } : {}
    }
    if (o.writes !== undefined) this.options.writes = strings(o.writes)
    if (o.origins !== undefined) this.options.origins = strings(o.origins)
    if (o.automations !== undefined) this.options.automations = strings(o.automations)
    if (o.chatbots !== undefined) this.options.chatbots = strings(o.chatbots)
    if (o.params !== undefined) {
      this.options.params =
        o.params && typeof o.params === 'object' && !Array.isArray(o.params) ? { ...o.params } : {}
    }
  }

  validate(): Array<BlockIssue> {
    const ee = super.validate()
    const app = this.options.applicationID && this.options.applicationID !== '0'
    if (!app && !this.options.source.trim()) {
      ee.push({ option: 'source', labelKey: 'block.issue.customSource' })
    }
    return ee
  }
}

Registry.set(kind, PageBlockCustom)
