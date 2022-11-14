/* eslint-disable @typescript-eslint/explicit-function-return-type */

import { Mark } from 'tiptap'
import { updateMark } from 'tiptap-commands'

/**
 * Represents text's color
 */
export default class Color extends Mark {
  get name () {
    return 'color'
  }

  get schema () {
    return {
      attrs: {
        color: {},
      },
      parseDOM: [
        {
          style: 'color',
          getAttrs: color => ({ color }),
        },
      ],

      toDOM: node => [
        'span',
        { style: `color: ${node.attrs.color}` },
        0,
      ],
    }
  }

  commands ({ type }) {
    return (attr) => updateMark(type, attr)
  }
}
