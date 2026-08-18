/**
 * The identity scopes the reference panel offers above the flow's own steps.
 *
 * The runtime injects them as fields of the global execution scope rather than
 * as step outputs, so `scope: "invoker"` is not a scope it can resolve — the
 * automation's own validator rejects it as `scope.unknown` and the TAQ never
 * registers. On the wire they are a dotted expression over the global scope.
 * The builder keeps addressing them as scopes, because that is how they are
 * grouped, listed and highlighted in the panel.
 */
const IDENTITY_SCOPES = ['invoker', 'runner']

const DOTTED = new RegExp(`^(${IDENTITY_SCOPES.join('|')})\\.(.+)$`)

interface ArgumentExpr {
  scope?: string
  expr?: string
  source?: string
  [k: string]: unknown
}

/** An argument as the builder holds it → as the API takes it. */
export function toWireArgument<T extends ArgumentExpr>(arg: T): T {
  if (!arg || !arg.scope || !IDENTITY_SCOPES.includes(arg.scope)) return arg

  const field = arg.expr || arg.source
  if (!field) return arg

  return { ...arg, scope: '', expr: `${arg.scope}.${field}`, source: undefined }
}

/** An argument as the API returns it → as the builder holds it. */
export function fromWireArgument<T extends ArgumentExpr>(arg: T): T {
  if (!arg || arg.scope) return arg

  const match = typeof arg.expr === 'string' ? arg.expr.match(DOTTED) : null
  if (!match) return arg

  return { ...arg, scope: match[1], expr: match[2] }
}

export const toWireArguments = <T extends ArgumentExpr>(args: T[] | undefined): T[] =>
  (args || []).map(toWireArgument)

export const fromWireArguments = <T extends ArgumentExpr>(args: T[] | undefined): T[] =>
  (args || []).map(fromWireArgument)
