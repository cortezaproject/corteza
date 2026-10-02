import { Apply, HumanID } from '../../../cast'

interface ButtonVisibility {
  // Server-evaluated expression; the button is shown only while it holds
  expression: string
}

export class Button {
  // Used when referring to Corredor automation script
  public script?: string = undefined

  // Used when referring to workflow with onManual trigger
  public workflowID?: string = undefined

  // Used when referring to NG automation
  public automationID?: string = undefined

  // Used when referring to a specific step (triggered by onManual trigger)
  public stepID?: string = undefined

  // Handle of the TAQ trigger this button was configured from
  public triggerHandle?: string = undefined

  // What the button runs: 'taq', 'workflow' or 'script'
  public scriptType?: string = undefined

  // resource type (copied from ui hook or from trigger)
  public resourceType?: string = undefined

  // Can override hook's label
  public label?: string = undefined

  // can override hook's variant
  public variant?: string = 'primary'

  public enabled = true

  // When the button is shown; absent or empty means always
  public visibility?: ButtonVisibility = undefined

  constructor(b: Partial<Button>) {
    Apply(this, b, Boolean, 'enabled')
    Apply(
      this,
      b,
      String,
      'label',
      'variant',
      'script',
      'resourceType',
      'triggerHandle',
      'scriptType',
    )
    Apply(this, b, HumanID, 'workflowID', 'stepID', 'automationID')

    if (b.visibility?.expression) {
      this.visibility = { expression: String(b.visibility.expression) }
    }
  }
}

export type PageBlockWrap = 'Plain' | 'Card'
