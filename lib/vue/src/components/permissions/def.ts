export const modalOpenEventName = 'c-permissions-modal-open'

interface ResourceParts {
  component: string;
  resourceType?: string;
  references: Array<string>;
  i18nPrefix: string;
}

/**
 * Splitting strings:
 *  - corteza::compose:moduleField/42/21/12
 *  - corteza::compose/42/21/12
 */
export function split (input: string): ResourceParts {
  const [tmp = '', references = ''] = input.split('/', 2)
  const [,, component, resourceType = undefined] = tmp.split(':')
  return {
    component,
    resourceType,
    references: references.split('/').filter(r => !!r),
    i18nPrefix: `resources.${component}.${resourceType || 'component'}`,
  }
}
