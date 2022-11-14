/* eslint-disable @typescript-eslint/explicit-function-return-type */

import { toggleBlockType } from 'tiptap-commands'
import { Paragraph as Base } from 'tiptap'
import { toAttrs, makeDOMParser } from './utils'

/**
 * Extends original paragraph node to allow content alignment
 */
export default class Paragraph extends Base {
  get schema () {
    return {
      ...super.schema,

      attrs: {
        ...super.schema.attrs,

        alignment: {
          default: undefined,
        },
      },

      parseDOM: super.schema.parseDOM.map(makeDOMParser),

      toDOM: (node) => [
        'p',
        toAttrs(node),
        0,
      ],
    }
  }

  commands ({ type }) {
    return (attrs) => toggleBlockType(type, type, attrs)
  }
}
