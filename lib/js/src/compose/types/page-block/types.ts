import { Apply, CortezaID } from '../../../cast'

interface ButtonVisibility {
  // Expression evaluated on the server; button is shown only when it is truthy
  expression: string;
}

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
  public variant?: string = 'primary';

  public enabled = true;

  // Controls when the button is shown
  public visibility?: ButtonVisibility = { expression: '' };

  constructor (b: Partial<Button>) {
    Apply(this, b, Boolean, 'enabled')
    Apply(this, b, String, 'label', 'variant', 'script', 'resourceType')
    Apply(this, b, CortezaID, 'workflowID', 'stepID')

    if (b.visibility) {
      this.visibility = { expression: String(b.visibility.expression || '') }
    }
  }
}

export type PageBlockWrap = 'Plain' | 'Card'
