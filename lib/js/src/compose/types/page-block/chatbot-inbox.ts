import { Apply } from '../../../cast'
import { PageBlock, PageBlockInput, Registry } from './base'

const kind = 'ChatbotInbox'

interface Options {
  chatbotIDs: string[]
  statusFilter: string[]
  refreshRate: number
  autoOpenFirst: boolean
}

const defaults: Readonly<Options> = Object.freeze({
  chatbotIDs: [],
  statusFilter: ['handoff_requested', 'handoff_active'],
  refreshRate: 5,
  autoOpenFirst: false,
})

export class PageBlockChatbotInbox extends PageBlock {
  readonly kind = kind

  options: Options = cloneDefaults()

  constructor(i?: PageBlockInput) {
    super(i)
    this.applyOptions(i?.options as Partial<Options>)
  }

  applyOptions(o?: Partial<Options>): void {
    if (!o) return

    if (Array.isArray(o.chatbotIDs)) {
      this.options.chatbotIDs = o.chatbotIDs.map(String)
    }
    if (Array.isArray(o.statusFilter) && o.statusFilter.length) {
      this.options.statusFilter = o.statusFilter.map(String)
    }

    Apply(this.options, o, Number, 'refreshRate')
    Apply(this.options, o, Boolean, 'autoOpenFirst')
  }
}

// cloneDefaults avoids sharing the frozen defaults object across block
// instances — apply() mutates options arrays, so each instance needs its own
// arrays to mutate.
function cloneDefaults(): Options {
  return {
    chatbotIDs: [...defaults.chatbotIDs],
    statusFilter: [...defaults.statusFilter],
    refreshRate: defaults.refreshRate,
    autoOpenFirst: defaults.autoOpenFirst,
  }
}

Registry.set(kind, PageBlockChatbotInbox)
