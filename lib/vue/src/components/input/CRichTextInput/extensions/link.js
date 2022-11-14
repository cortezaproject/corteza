/* eslint-disable @typescript-eslint/explicit-function-return-type */

import { Link as BaseLink } from 'tiptap-extensions'

/**
 * Extends original Link node to allow custom target attr
 */
export default class Link extends BaseLink {
  get schema () {
    const base = super.schema
    return {
      ...base,
      toDOM: node => ['a', {
        target: '_blank',
        ...node.attrs,
        rel: 'noopener noreferrer nofollow',
      }, 0],
    }
  }
}
