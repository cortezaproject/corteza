import { PageBlock, PageBlockInput, Registry, BlockIssue } from './base'
import { Button } from './types'
import { Apply } from '../../../cast'

const kind = 'Automation'

interface Options {
  // Ordered list of buttons to display in the block
  buttons: Array<Button>

  // When true, new compatible buttons (ui-hooks) are NOT
  // added automatically to the block
  //
  // Default behaviour is to add new buttons automatically.
  sealed: boolean
  magnifyOption: string
}

const defaults: Readonly<Options> = Object.freeze({
  buttons: [],
  sealed: false,
  magnifyOption: '',
})

export class PageBlockAutomation extends PageBlock {
  readonly kind = kind

  options: Options = { ...defaults }

  constructor(i?: PageBlockInput) {
    super(i)
    this.applyOptions(i?.options as Partial<Options>)
  }

  applyOptions(o?: Partial<Options>): void {
    if (!o) return

    Apply(this.options, o, String, 'magnifyOption')

    if (o.buttons) {
      this.options.buttons = o.buttons.map(b => new Button(b))
    }
  }

  // Validates Page Block configuration
  validate(): Array<BlockIssue> {
    const ee = super.validate()

    if (!this.options.buttons.length) {
      ee.push({ option: 'buttons', labelKey: 'block.issue.noAutomationButtons' })
    }

    this.options.buttons.forEach((b, i) => {
      if (b.workflowID) {
        // workflow defined
        return
      }

      if (b.automationID) {
        // NG automation defined
        return
      }

      if (b.script) {
        // script defined
        return
      }

      ee.push({
        option: `buttons.${i}`,
        labelKey: 'block.issue.unconfiguredAutomationButton',
      })
    })

    return ee
  }
}

Registry.set(kind, PageBlockAutomation)
