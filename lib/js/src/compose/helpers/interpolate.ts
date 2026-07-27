export interface InterpolationVariables {
  record?: unknown
  user?: unknown
  recordID?: string | number
  ownerID?: string | number
  userID?: string | number
}

// Evaluates a JS template literal string (e.g. `id = ${recordID}`) against
// compose's record/user variables. `template` is trusted page-author
// configuration — a prefilter, URL, label, etc. typed into the page/field
// builder by an admin, never end-user input — which is why building a
// function from it is acceptable here.
//
// The variables are passed as explicit parameters rather than picked up from
// the enclosing scope: bindings referenced only from inside the template are
// invisible to bundlers and may be dropped or renamed when minified.
export function interpolateTemplate(
  template: string,
  { record, user, recordID, ownerID, userID }: InterpolationVariables,
): string {
  const evaluate = new Function(
    'record',
    'user',
    'recordID',
    'ownerID',
    'userID',
    'return `' + template + '`',
  )

  return evaluate(record, user, recordID, ownerID, userID)
}
