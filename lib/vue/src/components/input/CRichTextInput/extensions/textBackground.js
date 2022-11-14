/* eslint-disable @typescript-eslint/explicit-function-return-type */

import { Mark } from 'tiptap'
import { updateMark } from 'tiptap-commands'

/**
 * Represents text's background color
 */
export default class Background extends Mark {
  get name () {
    return 'background'
  }

  get schema () {
    return {
      attrs: {
        color: {},
      },

      parseDOM: [
        {
          style: 'background-color',
          getAttrs: color => ({ color }),
        },
      ],

      toDOM: node => [
        'span',
        { style: `background-color: ${node.attrs.color}` },
        0,
      ],
    }
  }

  commands ({ type }) {
    return (attr) => updateMark(type, attr)
  }
}
