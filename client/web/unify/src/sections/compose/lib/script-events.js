/**
 * Matches one trigger constraint of a Corredor script against the namespace and
 * the module the script is dispatched for.
 *
 * The event carries both because the module a record is built from does not know
 * the namespace it lives in: modules are cached per namespace and frozen, so the
 * namespace cannot be read off them and cannot be set on them either.
 *
 * Keep the property names in sync with server/compose/service/event.
 */
export function scriptConstraintMatcher({ namespace, module } = {}) {
  return c => {
    switch (c.Name()) {
      case 'namespace':
      case 'namespace.slug':
        return c.Match(namespace?.slug || '')
      case 'namespace.name':
        return c.Match(namespace?.name || '')
      case 'module':
      case 'module.handle':
        return c.Match(module?.handle || '')
      case 'module.name':
        return c.Match(module?.name || '')
      default:
        return false
    }
  }
}

/**
 * Whether a rejected script dispatch is a script refusing the action
 */
export function isScriptAbort(e) {
  const message = e?.message || ''
  return message.toLowerCase() === 'aborted'
}
