export interface ScriptConstraint {
  name?: string
  op?: string
  value?: string[]
}

// The operator a constraint carries when the script named none
const defaultConstraintOp = '='

/**
 * One chip per trigger constraint, as `<name> <op> <value>`; an op the script
 * left out is the equality the server assumes.
 */
export function constraintChips(constraints?: ScriptConstraint[] | null): string[] {
  return (constraints || [])
    .map(constraint =>
      [
        constraint?.name,
        constraint?.op || defaultConstraintOp,
        (constraint?.value || []).join(', '),
      ]
        .filter(Boolean)
        .join(' '),
    )
    .filter(chip => chip !== defaultConstraintOp)
}
