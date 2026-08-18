import { Apply } from '../../../cast'
import { PageBlock, PageBlockInput, Registry, BlockIssue } from './base'

const kind = 'AgentChat'

interface Options {
  allowedAgentIDs: string[]
  defaultAgentID: string
  autoResume: boolean
}

const defaults: Readonly<Options> = Object.freeze({
  allowedAgentIDs: [],
  defaultAgentID: '',
  autoResume: true,
})

export class PageBlockAgentChat extends PageBlock {
  readonly kind = kind

  options: Options = cloneDefaults()

  constructor(i?: PageBlockInput) {
    super(i)
    this.applyOptions(i?.options as Partial<Options>)
  }

  applyOptions(o?: Partial<Options>): void {
    if (!o) return

    if (Array.isArray(o.allowedAgentIDs)) {
      this.options.allowedAgentIDs = o.allowedAgentIDs.map(String)
    }

    Apply(this.options, o, String, 'defaultAgentID')
    Apply(this.options, o, Boolean, 'autoResume')
  }

  validate(): Array<BlockIssue> {
    const ee = super.validate()
    if (!this.options.allowedAgentIDs.length) {
      ee.push({ option: 'allowedAgentIDs', labelKey: 'block.issue.noAgents' })
    }
    return ee
  }
}

function cloneDefaults(): Options {
  return {
    allowedAgentIDs: [...defaults.allowedAgentIDs],
    defaultAgentID: defaults.defaultAgentID,
    autoResume: defaults.autoResume,
  }
}

Registry.set(kind, PageBlockAgentChat)
