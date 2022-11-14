/* eslint-disable @typescript-eslint/explicit-function-return-type */
import { toggleBlockType } from 'tiptap-commands'
import { Heading as Base } from 'tiptap-extensions'
import { toAttrs, makeDOMParser } from './utils'

/**
 * Extends original heading node to allow content alignment
 */
export default class Heading extends Base {
  get defaultOptions () {
    return {
      ...super.defaultOptions,
      alignment: undefined,
    }
  }

  get baseP () {
    return this.options.levels
      .map(level => ({
        tag: `h${level}`,
        getAttrs: () => ({ level }),
      }))
  }

  get schema () {
    return {
      ...super.schema,

      attrs: {
        ...super.schema.attrs,

        alignment: {
          default: undefined,
        },
      },

      parseDOM: this.baseP.map(makeDOMParser),

      toDOM: (node) => [
        `h${node.attrs.level}`,
        toAttrs(node),
        0,
      ],
    }
  }

  commands ({ type, schema }) {
    return (attrs) => {
      return toggleBlockType(type, schema.nodes.paragraph, attrs)
    }
  }
}
