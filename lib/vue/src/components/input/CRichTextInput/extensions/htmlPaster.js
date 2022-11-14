/* eslint-disable @typescript-eslint/explicit-function-return-type */

import hps from 'html-parse-stringify'
import {
  Extension,
  Plugin,
  PluginKey,
} from 'tiptap'

const knownContentRoots = new Set(['pre', 'p'])

function tp (c) {
  if (c.type && c.type.name === 'text') {
    c.text = (c.text || '').replace(/\u00a0/ig, '\u0020')
  }

  for (const cc of c.content.content) {
    tp(cc)
  }

  return c
}

/**
 * This plugin handles HTML pasting.
 * It prevents white spaces from collapsing by converting them into
 * non-breakable spaces.
 */
const plug = new Plugin({
  key: new PluginKey('htmlPaster'),

  view () {
    return {
      // eslint-disable-next-line @typescript-eslint/no-empty-function
      update: () => {},
    }
  },

  state: {
    init () {
      return {}
    },

    apply (tr, prev) {
      return { ...prev }
    },
  },

  props: {
    /**
     * Method is invoked when a user pastes HTML content into the editor.
     * The method replaces regular white spaces (\u0020) with non-breakable white spaces (\u00a0).
     * This forces the underlying HTML parser to not collapse white spaces.
     *
     * It also performs a slight HTML node manipulation to assure consistent behaviour between different
     * environments.
     *
     * @param {String} html Pasted HTML content
     */
    transformPastedHTML (html) {
      const tree = hps.parse(html).filter(({ type }) => type !== 'text')
      if (tree.length === 1) {
        // 1. find main content node
        let root
        if (tree[0].name === 'html') {
          root = tree[0].children.find(({ name }) => name === 'body' || name === 'main')
        } else {
          root = tree[0]
        }

        if (root) {
          // 2. make content nodes consistent
          root.children.forEach((c, i) => {
            if (!c.voidElement && !knownContentRoots.has(c.name)) {
              c.name = 'p'
            }
            if (c.name === 'br') {
              root.children[i].voidElement = false
              root.children[i].name = 'p'
            }
          })

          html = hps.stringify(tree)
        }
      }

      html = html.replace(/&#32;/ig, '\u00a0')
      html = html.replace(/\n/ig, '<br>')
      // lets be smart and only replace trailing/leading white spaces.
      // replacing every white space can mess with other plugins (hint hint mentions)
      const rr = /(\u0020+)<|>(\u0020+)/gi
      let exec
      do {
        exec = rr.exec(html)
        if (exec) {
          const cc = exec[1] !== undefined ? exec[1] : exec[2]
          html = html.substring(0, exec.index + 1) +
            new Array(cc.length + 1).join('\u00a0') +
            html.substring(exec.index + cc.length + (exec[1] === undefined ? 1 : 0))
        }
      } while (exec)

      return html
    },

    /**
     * Method is invoked before the given content slice is inserted into the document.
     * The method replaces the above created non-breakable white spaces (\u00a0) with regular white spaces (\u0020).
     *
     * @param {Slice} c Pasted content slice
     */
    transformPasted (c) {
      return tp(c)
    },
  },
})

export default class HtmlPaster extends Extension {
  get name () {
    return 'html-paster'
  }

  get plugins () {
    return [plug]
  }
}
