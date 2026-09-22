import { eventbus } from '@planetcrust/human-js'

/**
 * Matches one trigger constraint of a Corredor script against the namespace and
 * the module the script is dispatched for.
 *
 * The event carries both rather than reading them off the record: a module knows
 * its namespace only when the namespace store held it at the time the module was
 * cached, and a page can dispatch before that.
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

/**
 * The compose resources a page can hand a manually triggered script.
 *
 * Every page carries itself, its module and its namespace; only a record page
 * carries a record, so a script bound to `compose:record` has nothing to run
 * against anywhere else.
 */
export function pageScriptResourceTypes(page) {
  const resourceTypes = ['compose', 'compose:namespace', 'compose:page', 'compose:module']

  if (page?.isRecordPage) {
    resourceTypes.push('compose:record')
  }

  return resourceTypes
}

/**
 * Whether a page can satisfy the resource a trigger is bound to.
 */
export function pageFitsResourceType(page, resourceType) {
  if (!resourceType) return true
  return pageScriptResourceTypes(page).includes(resourceType)
}

// The constraint names a page can decide on, and the value it holds for each.
function constraintContext({ namespace, module } = {}) {
  return {
    namespace: namespace?.slug,
    'namespace.slug': namespace?.slug,
    'namespace.name': namespace?.name,
    module: module?.handle,
    'module.handle': module?.handle,
    'module.name': module?.name,
  }
}

/**
 * Whether a trigger's constraints leave it any chance of firing on this page.
 *
 * Only a constraint the page can settle counts against it: one naming a
 * namespace or module property the page knows the value of, and not matching it.
 * An unknown property, a value the page does not hold and a constraint that
 * cannot be read all leave the trigger possible — unlike dispatch, where an
 * unknown value is matched as empty and the script is simply not run.
 */
export function triggerCanApply(trigger, { namespace, module } = {}) {
  const context = constraintContext({ namespace, module })

  return (trigger?.constraints || []).every(constraint => {
    let matcher

    try {
      matcher = eventbus.ConstraintMaker(constraint)
    } catch {
      return true
    }

    const name = matcher.Name()

    if (!(name in context) || context[name] === undefined || context[name] === null) {
      return true
    }

    return matcher.Match(context[name])
  })
}
