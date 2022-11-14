import { Apply, CortezaID } from '../../../cast'

export class Button {
  // Used when referring to Corredor automation script
  public script?: string = undefined

  // Used when referring to workflow with onManual trigger
  public workflowID?: string = undefined

  // Used when referring to a specific step (triggered by onManual trigger)
  public stepID?: string = undefined

  // resource type (copied from ui hook or from trigger)
  public resourceType?: string = undefined;

  // Can override hook's label
  public label?: string = undefined;

  // can override hook's variant
  public variant?: string = undefined;

  public enabled = true;

  constructor (b: Partial<Button>) {
    Apply(this, b, Boolean, 'enabled')
    Apply(this, b, String, 'label', 'variant', 'script', 'resourceType')
    Apply(this, b, CortezaID, 'workflowID', 'stepID')
  }
}

export type PageBlockWrap = 'Plain' | 'Card'
