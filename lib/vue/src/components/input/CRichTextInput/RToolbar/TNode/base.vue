<script>
import { nodeTypes } from '../../lib/formats'

/**
 * Defines common props, methods between different toolbar item types
 */
export default {
  props: {
    format: {
      type: Object,
      required: true,
      default: () => ({}),
    },
  },

  methods: {
    /**
     * Helper method to emit format selection
     * @param {String} type Format's type
     * @param {Object} attrs Format's extra attributes
     */
    onClick (type, attrs) {
      if (!attrs) {
        attrs = {}
      }

      // manual toggling, since tiptap's isActive fails to resolve modified nodes
      const act = this.activeNode([type], attrs)
      if (act) {
        type = 'paragraph'
      }

      const nn = this.activeNode(nodeTypes)
      if (!nn) {
        throw new Error('no node selected')
      }
      const n = nn.node

      // preserve some attrs
      const cAttr = n.attrs
      const target = this.$attrs.editor.nodes[type]
      const nAttrs = {}

      for (const a in target.attrs) {
        if (attrs[a] === undefined) {
          nAttrs[a] = cAttr[a]
        } else {
          nAttrs[a] = attrs[a]
        }
      }

      this.$emit('click', { type, attrs: nAttrs })
    },

    /**
     * Helper to determine active nodes in the given selection.
     * Replaces tiptap's built in isActive, since it can't handle our modified
     * nodes.
     */
    activeNodes (types, attrs) {
      const ed = this.$attrs.editor
      const rtr = []
      ed.state.doc.nodesBetween(
        ed.selection.from,
        ed.selection.to,
        (n, pos) => {
          if (types.includes(n.type.name)) {
            if (attrs) {
              if (!Object.entries(attrs || {}).find(([k, v]) => n.attrs[k] !== v)) {
                rtr.push({ node: n, position: pos })
              }
            } else {
              rtr.push({ node: n, position: pos })
            }
          }
        },
      )

      return rtr
    },

    /**
     * Helper to determine active node in the given selection
     */
    activeNode (types, attrs) {
      const ann = this.activeNodes(types, attrs)
      if (!ann) {
        return undefined
      }
      return ann[0]
    },

    isActiveCheck (types, attrs) {
      return !!this.activeNode(types, attrs)
    },

    /**
     * Helper method to determine if given format is active or not.
     * When attrs is provided, it will check for that exact match
     * @param {Object|undefined} attrs Format's extra attributes
     * @returns {Array|undefined}
     */
    activeClasses (attrs) {
      if (this.isActiveCheck([this.format.type], attrs)) {
        return ['text-success']
      }

      return undefined
    },
  },
}
</script>
